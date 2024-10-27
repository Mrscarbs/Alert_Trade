use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use env_logger::Env;
use log::{error, info};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use sqlx::{Pool, MySql, Row};
use sqlx::mysql::MySqlPoolOptions;
use sqlx::Error;
use std::fs::OpenOptions;
use std::io::Write;
use chrono::Utc;

// Request struct for /get-history-bars API with interval
#[derive(Serialize, Deserialize, Debug)]
struct HistoryBarsRequest {
    symbol: String,
    bidask: u8,
    from: String,       // format: yymmddTHH:mm:ss
    to: String,         // format: yymmddTHH:mm:ss
    interval: String,   // e.g., "1min", "5min", "1hour"
    response_type: String, // either "csv" or "json"
}

// Request struct for /get-top-gainers and /get-top-losers API
#[derive(Serialize, Deserialize, Debug)]
struct TopNRecordsRequest {
    segment: String,
    response_type: String,
    topn: u32,
}

// Updated TopRecord struct with manual serialization
#[derive(Debug, Serialize)]
struct TopRecord {
    name: String,
    id: u64,
    date: String,
    open: f64,
    high: f64,
    low: f64,
    close: f64,
    volume: u64,
    value: f64,
    prev_close: f64,
    change: f64,
    percent_change: f64,
    series: String,
    exchange: String,
}

// Updated response struct for array-based responses
#[derive(Debug, Serialize, Deserialize)]
struct TrueDataTopRecordsResponse {
    status: String,
    #[serde(rename = "Records")]
    records: Vec<Vec<Value>>,
}

// Function to convert array response to TopRecord
impl TopRecord {
    fn from_array(arr: &[Value]) -> Option<Self> {
        if arr.len() < 14 {
            return None;
        }

        Some(TopRecord {
            name: arr[0].as_str()?.to_string(),
            id: arr[1].as_u64()?,
            date: arr[2].as_str()?.to_string(),
            open: arr[3].as_f64()?,
            high: arr[4].as_f64()?,
            low: arr[5].as_f64()?,
            close: arr[6].as_f64()?,
            volume: arr[7].as_u64()?,
            value: arr[8].as_f64()?,
            prev_close: arr[9].as_f64()?,
            change: arr[10].as_f64()?,
            percent_change: arr[11].as_f64()?,
            series: arr[12].as_str()?.to_string(),
            exchange: arr[13].as_str()?.to_string(),
        })
    }
}

// Function to fetch historical bars from TrueData with interval
async fn fetch_history_bars(client: &Client, params: &HistoryBarsRequest, access_token: &str) -> Result<String, reqwest::Error> {
    let url = format!(
        "https://history.truedata.in/getbars?symbol={}&bidask={}&from={}&to={}&interval={}&response={}",
        params.symbol, params.bidask, params.from, params.to, params.interval, params.response_type
    );

    let response = client
        .get(&url)
        .bearer_auth(access_token)
        .send()
        .await?
        .text()
        .await?;

    Ok(response)
}

// Updated function to fetch top N gainers from TrueData
async fn fetch_top_n_gainers(client: &Client, params: &TopNRecordsRequest, access_token: &str) -> Result<TrueDataTopRecordsResponse, reqwest::Error> {
    let url = format!(
        "https://history.truedata.in/gettopngainers?segment={}&response={}&topn={}",
        params.segment, params.response_type, params.topn
    );

    let response = client
        .get(&url)
        .bearer_auth(access_token)
        .header("accept", "application/json")
        .send()
        .await?
        .json::<TrueDataTopRecordsResponse>()
        .await?;

    Ok(response)
}

// Updated function to fetch top N losers from TrueData
async fn fetch_top_n_losers(client: &Client, params: &TopNRecordsRequest, access_token: &str) -> Result<TrueDataTopRecordsResponse, reqwest::Error> {
    let url = format!(
        "https://history.truedata.in/gettopnlosers?segment={}&response={}&topn={}",
        params.segment, params.response_type, params.topn
    );

    let response = client
        .get(&url)
        .bearer_auth(access_token)
        .header("accept", "application/json")
        .send()
        .await?
        .json::<TrueDataTopRecordsResponse>()
        .await?;

    Ok(response)
}

// Function to connect to MySQL database
async fn connect_to_mysql_db() -> Result<Pool<MySql>, Error> {
    let database_url = "mysql://root:Karma100%25@alerttrade.cbgqgqswkxrn.eu-north-1.rds.amazonaws.com:3306/alert_trade_db";
    let pool = MySqlPoolOptions::new()
        .max_connections(5)
        .connect(database_url)
        .await?;
    Ok(pool)
}

// Function to get access token from MySQL
async fn get_access_token(pool: &Pool<MySql>, api_id: i64) -> Result<String, Error> {
    let row = sqlx::query("CALL stp_get_access_token_api_id(?)")
        .bind(api_id)
        .fetch_one(pool)
        .await?;

    let access_token: String = row.try_get(0)?;
    Ok(access_token)
}

// Function to log requests and errors to a file
fn log_to_file(message: &str) {
    let mut log_file = OpenOptions::new()
        .append(true)
        .create(true)
        .open("request_logs.log")
        .unwrap();

    let timestamp = Utc::now();
    writeln!(log_file, "{} - {}", timestamp, message).unwrap();
}

// Handler for /get-history-bars API
async fn handle_get_history_bars(
    pool: web::Data<Pool<MySql>>,
    params: web::Query<HistoryBarsRequest>,
) -> impl Responder {
    let client = Client::new();

    let log_message = format!("Received /get-history-bars request: {:?}", params);
    log_to_file(&log_message);
    info!("{}", log_message);

    let access_token = match get_access_token(pool.get_ref(), 1).await {
        Ok(token) => token,
        Err(e) => {
            let error_message = format!("Failed to get access token: {}", e);
            error!("{}", error_message);
            log_to_file(&error_message);
            return HttpResponse::InternalServerError().body("Failed to get access token");
        }
    };

    match fetch_history_bars(&client, &params, &access_token).await {
        Ok(response) => {
            let success_message = "Successfully fetched data from TrueData API".to_string();
            log_to_file(&success_message);
            info!("{}", success_message);
            HttpResponse::Ok().body(response)
        },
        Err(e) => {
            let error_message = format!("Failed to fetch data from TrueData API: {}", e);
            error!("{}", error_message);
            log_to_file(&error_message);
            HttpResponse::InternalServerError().body("Failed to fetch data from TrueData API")
        }
    }
}

// Updated handler for /get-top-gainers API
async fn handle_get_top_n_gainers(
    pool: web::Data<Pool<MySql>>,
    params: web::Query<TopNRecordsRequest>,
) -> impl Responder {
    let client = Client::new();

    let log_message = format!("Received /get-top-gainers request: {:?}", params);
    log_to_file(&log_message);
    info!("{}", log_message);

    let access_token = match get_access_token(pool.get_ref(), 1).await {
        Ok(token) => token,
        Err(e) => {
            let error_message = format!("Failed to get access token: {}", e);
            error!("{}", error_message);
            log_to_file(&error_message);
            return HttpResponse::InternalServerError().body("Failed to get access token");
        }
    };

    match fetch_top_n_gainers(&client, &params, &access_token).await {
        Ok(response) => {
            let records: Vec<TopRecord> = response.records
                .iter()
                .filter_map(|arr| TopRecord::from_array(arr))
                .collect();

            let success_message = "Successfully fetched top N gainers from TrueData API".to_string();
            log_to_file(&success_message);
            info!("{}", success_message);
            
            let formatted_response = serde_json::json!({
                "status": response.status,
                "Records": records
            });
            
            HttpResponse::Ok().json(formatted_response)
        },
        Err(e) => {
            let error_message = format!("Failed to fetch top N gainers from TrueData API: {}", e);
            error!("{}", error_message);
            log_to_file(&error_message);
            HttpResponse::InternalServerError().body("Failed to fetch data from TrueData API")
        }
    }
}

// Updated handler for /get-top-losers API
async fn handle_get_top_n_losers(
    pool: web::Data<Pool<MySql>>,
    params: web::Query<TopNRecordsRequest>,
) -> impl Responder {
    let client = Client::new();

    let log_message = format!("Received /get-top-losers request: {:?}", params);
    log_to_file(&log_message);
    info!("{}", log_message);

    let access_token = match get_access_token(pool.get_ref(), 1).await {
        Ok(token) => token,
        Err(e) => {
            let error_message = format!("Failed to get access token: {}", e);
            error!("{}", error_message);
            log_to_file(&error_message);
            return HttpResponse::InternalServerError().body("Failed to get access token");
        }
    };

    match fetch_top_n_losers(&client, &params, &access_token).await {
        Ok(response) => {
            let records: Vec<TopRecord> = response.records
                .iter()
                .filter_map(|arr| TopRecord::from_array(arr))
                .collect();

            let success_message = "Successfully fetched top N losers from TrueData API".to_string();
            log_to_file(&success_message);
            info!("{}", success_message);
            
            let formatted_response = serde_json::json!({
                "status": response.status,
                "Records": records
            });
            
            HttpResponse::Ok().json(formatted_response)
        },
        Err(e) => {
            let error_message = format!("Failed to fetch top N losers from TrueData API: {}", e);
            error!("{}", error_message);
            log_to_file(&error_message);
            HttpResponse::InternalServerError().body("Failed to fetch data from TrueData API")
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize the logger
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();

    // Initialize the MySQL connection pool
    let pool = connect_to_mysql_db().await.expect("Failed to create MySQL pool");

    // Log the start of the server
    let startup_message = "Starting the Actix Web server on http://0.0.0.0:8080".to_string();
    log_to_file(&startup_message);
    info!("{}", startup_message);

    // Start Actix Web server with all three routes
    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(pool.clone()))
            .route("/get-history-bars", web::get().to(handle_get_history_bars))
            .route("/get-top-gainers", web::get().to(handle_get_top_n_gainers))
            .route("/get-top-losers", web::get().to(handle_get_top_n_losers))
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}