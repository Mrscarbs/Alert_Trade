use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use sqlx::{Pool, MySql, Row};
use sqlx::mysql::MySqlPoolOptions;

#[derive(Deserialize)]
struct EntityId {
    entity_id: i64,
}

#[derive(Serialize)]
struct ShareDetails {
    candidate: String,
    share_name: String,
    amount: i32,
}

// Modified database connection function with proper URL format
async fn connect_to_mysql_db() -> Result<Pool<MySql>, sqlx::Error> {
    // Corrected database URL format
    let database_url = "mysql://root:Karma100%25@alerttrade.cbgqgqswkxrn.eu-north-1.rds.amazonaws.com:3306/alert_trade_db";
    
    MySqlPoolOptions::new()
        .max_connections(5)
        .connect(database_url)
        .await
}

async fn execute_stored_procedure(pool: &Pool<MySql>, entity_id: i64) -> Result<Vec<ShareDetails>, sqlx::Error> {
    let rows = sqlx::query(
        "CALL stp_get_candidate_share_details(?)"
    )
    .bind(entity_id)
    .fetch_all(pool)
    .await?;

    let mut share_details = Vec::new();

    for row in rows {
        share_details.push(ShareDetails {
            candidate: row.try_get(0)?,
            share_name: row.try_get(1)?,
            amount: row.try_get::<i32, _>(2)?,
        });
    }

    Ok(share_details)
}

async fn get_candidate_share_details(data: web::Json<EntityId>, pool: web::Data<Pool<MySql>>) -> impl Responder {
    let entity_id = data.entity_id;
    match execute_stored_procedure(pool.get_ref(), entity_id).await {
        Ok(details) => HttpResponse::Ok().json(details),
        Err(e) => {
            eprintln!("Database error: {:?}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to fetch data from database",
                "details": e.to_string()
            }))
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Enhanced error handling for database connection
    let pool = match connect_to_mysql_db().await {
        Ok(pool) => {
            println!("Successfully connected to database");
            pool
        },
        Err(e) => {
            eprintln!("Failed to connect to database: {:?}", e);
            std::process::exit(1);
        }
    };
    
    let pool_data = web::Data::new(pool);
    
    println!("Starting server on 0.0.0.0:8080");
    
    HttpServer::new(move || {
        App::new()
            .app_data(pool_data.clone())
            .route("/get_share_details", web::post().to(get_candidate_share_details))
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}