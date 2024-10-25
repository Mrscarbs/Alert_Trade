use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use sqlx::{Pool, MySql, Row}; // Use Row to fetch fields by index
use sqlx::mysql::MySqlPoolOptions;

// Struct to represent the input entity ID
#[derive(Deserialize)]
struct EntityId {
    entity_id: i64,
}

// Struct for the candidate share details
#[derive(Serialize)]
struct ShareDetails {
    candidate: String,
    share_name: String,
    amount: i32,  // Assuming 'amount' is an integer in your database
}

// Function to connect to the MySQL database with a hardcoded URL
async fn connect_to_mysql_db() -> Pool<MySql> {
    // Hardcoded database URL
    let database_url = "mysql://root:Karma100%@tcp(alerttrade.cbgqgqswkxrn.eu-north-1.rds.amazonaws.com:3306)/alert_trade_db";
    
    MySqlPoolOptions::new()
        .max_connections(5)
        .connect(database_url)
        .await
        .expect("Failed to connect to the database")
}

// Function to call the stored procedure and get share details (without compile-time validation)
async fn execute_stored_procedure(pool: &Pool<MySql>, entity_id: i64) -> Result<Vec<ShareDetails>, sqlx::Error> {
    let rows = sqlx::query(
        "CALL stp_get_candidate_share_details(?)" // Using query() instead of query!() to avoid compile-time validation
    )
    .bind(entity_id)
    .fetch_all(pool)
    .await?;

    let mut share_details = Vec::new();

    // Iterate through the rows and access the fields by their index
    for row in rows {
        share_details.push(ShareDetails {
            candidate: row.try_get(0)?,           // Assuming 'candidate' is at index 0 (String)
            share_name: row.try_get(1)?,          // Assuming 'share_name' is at index 1 (String)
            amount: row.try_get::<i32, _>(2)?,    // Fetch 'amount' as an integer (index 2)
        });
    }

    Ok(share_details)
}

// API handler to process the request
async fn get_candidate_share_details(data: web::Json<EntityId>, pool: web::Data<Pool<MySql>>) -> impl Responder {
    let entity_id = data.entity_id;
    let result = execute_stored_procedure(pool.get_ref(), entity_id).await;
    
    match result {
        Ok(details) => HttpResponse::Ok().json(details),
        Err(_) => HttpResponse::InternalServerError().body("Error fetching data"),
    }
}

// Main function to start the server
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let pool = connect_to_mysql_db().await;
    let pool_data = web::Data::new(pool); // Share the DB pool between routes
    
    HttpServer::new(move || {
        App::new()
            .app_data(pool_data.clone())
            .route("/get_share_details", web::post().to(get_candidate_share_details))
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}