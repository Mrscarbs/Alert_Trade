use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use log::{info, error, debug}; // Import logging macros
use sqlx::mysql::MySqlPoolOptions;
use sqlx::{MySql, Pool, Row};  // Import the Row trait
use sqlx::Error;
use actix_web::{post, web, App, HttpServer, HttpResponse};
use env_logger::Env;  // Import env_logger for initializing logger

#[derive(Serialize, Clone)]
struct RequestBody {
    list_symbols: Vec<String>,
    entity_id: i64,
}

#[derive(Debug, Deserialize, Clone)]
struct PriceDetails {
    Closes: f64,
    High: f64,
    Low: f64,
    Timestamp: i64,
    Position_price: f64,
    Position_quantity: i64,
}

#[derive(Debug, Deserialize, Clone)]
struct CompanyData {
    Equity: Vec<f64>,
    Trade_details: Vec<PriceDetails>,
}

#[derive(Debug, Deserialize, Clone)]
struct ApiResponse(HashMap<String, CompanyData>);

#[derive(Debug, Serialize, Clone)]
struct StockAnalysis {
    symbol: String,
    sharpe_ratio: f64,
    sortino_ratio: f64,
    Timestamp: i64,
    Equity: f64,
}

#[derive(Debug, Serialize, Clone)]
struct PortfolioMetrics {
    cumulative_analysis: Vec<StockAnalysis>,
    top_gainers_losers: Vec<(String, f64)>,
    portfolio_allocation: Vec<(String, f64, f64)>,
    total_equity_over_time: Vec<f64>,
    total_sharpe_ratio: HashMap<i64, f64>,
    total_sortino_ratio: HashMap<i64, f64>,
}

#[derive(Deserialize)]
struct InputData {
    entity_id: i64,
}

#[post("/user_analytics_info")]
async fn user_analytics_info(input: web::Json<InputData>) -> HttpResponse {
    let entity_id = input.entity_id;  // Extract entity_id before moving input
    info!("Received request with entity_id: {}", entity_id);

    match main_logic(input.into_inner()).await {
        Ok(result) => {
            info!("Successfully processed request for entity_id: {}", entity_id);
            HttpResponse::Ok().json(result)
        },
        Err(e) => {
            error!("Failed to process request for entity_id: {}: {}", entity_id, e);
            HttpResponse::InternalServerError().body(e.to_string())
        },
    }
}

async fn main_logic(input: InputData) -> Result<PortfolioMetrics, Box<dyn std::error::Error>> {
    let client = Client::new();
    let pool = connect_to_db().await?;
    let list_symbols = execute_stored_procedure(&pool, input.entity_id).await?;

    let body = RequestBody {
        list_symbols,
        entity_id: input.entity_id,
    };

    info!("Sending request to external API for entity_id: {}", input.entity_id);

    let response = client
        .post("http://localhost:8080/get_historical_performance")
        .json(&body)
        .send()
        .await?;

    let response_text = response.text().await?;
    let trimmed_response_text = response_text.trim();

    let api_response: ApiResponse = serde_json::from_str(&trimmed_response_text)
        .map_err(|e| {
            error!("Failed to parse API response for entity_id: {}: {}", input.entity_id, e);
            Box::new(e) as Box<dyn std::error::Error>
        })?;

    let mut performances = Vec::new();

    for (symbol, data) in &api_response.0 {
        let prices: Vec<f64> = data.Trade_details.iter().map(|pd| if pd.Closes != 0.0 { pd.Closes } else { pd.Position_price }).collect();
        let Timestamps: Vec<i64> = data.Trade_details.iter().map(|pd| pd.Timestamp).collect();
        let Equity: Vec<f64> = data.Equity.clone();

        debug!("Calculating returns for symbol in entity_id: {}", input.entity_id);
        let returns = calculate_returns(&prices);
        let analysis = cumulative_analysis(&returns, &Timestamps, &Equity, symbol);
        performances.extend(analysis);
    }

    performances.sort_by(|a, b| a.Timestamp.cmp(&b.Timestamp));

    info!("Creating timestamp list for entity_id: {}", input.entity_id);
    let vec_Timestamp_sorted = all_Timestamp_list_creator(&api_response);

    info!("Creating analysis map for entity_id: {}", input.entity_id);
    let analysis_map = create_analysis_map(&vec_Timestamp_sorted, &performances);

    info!("Calculating portfolio metrics for entity_id: {}", input.entity_id);
    let top_gainers_losers = get_top_gainers_losers(&api_response, 5);
    let portfolio_allocation = get_portfolio_allocation(&api_response);
    let total_equity_over_time = get_total_Equity(&api_response);
    let total_sharpe_ratio = calculate_total_sharpe_ratio(&analysis_map);
    let total_sortino_ratio = calculate_total_sortino_ratio(&analysis_map);

    let portfolio_metrics = PortfolioMetrics {
        cumulative_analysis: performances,
        top_gainers_losers,
        portfolio_allocation,
        total_equity_over_time,
        total_sharpe_ratio,
        total_sortino_ratio,
    };

    info!("Successfully calculated metrics for entity_id: {}", input.entity_id);

    Ok(portfolio_metrics)
}

fn calculate_returns(prices: &[f64]) -> Vec<f64> {
    prices.windows(2)
        .map(|w| (w[1] - w[0]) / w[0])
        .collect()
}

fn cumulative_analysis(returns: &[f64], Timestamps: &[i64], Equity: &[f64], symbol: &str) -> Vec<StockAnalysis> {
    (2..=returns.len())
        .map(|i| {
            let sharpe_ratio = calculate_sharpe_ratio(&returns[..i]);
            let sortino_ratio = calculate_sortino_ratio(&returns[..i]);
            StockAnalysis {
                symbol: symbol.to_string(),
                sharpe_ratio,
                sortino_ratio,
                Timestamp: Timestamps[i],
                Equity: Equity[i],
            }
        })
        .collect()
}

fn calculate_sharpe_ratio(returns: &[f64]) -> f64 {
    let mean: f64 = returns.iter().sum::<f64>() / returns.len() as f64;
    let std_dev: f64 = (returns.iter().map(|r| (r - mean).powi(2)).sum::<f64>() / returns.len() as f64).sqrt();
    if std_dev == 0.0 {
        return 0.0;
    }
    mean / std_dev
}

fn calculate_sortino_ratio(returns: &[f64]) -> f64 {
    let mean: f64 = returns.iter().sum::<f64>() / returns.len() as f64;
    let negative_std_dev: f64 = (returns.iter().filter(|r| **r < 0.0).map(|r| r.powi(2)).sum::<f64>() / returns.len() as f64).sqrt();
    if negative_std_dev == 0.0 {
        return 0.0;
    }
    mean / negative_std_dev
}

fn all_Timestamp_list_creator(stock_map: &ApiResponse) -> Vec<i64> {
    let mut temp_vec = Vec::new();
    for (_key, value) in &stock_map.0 {
        temp_vec.extend(value.Trade_details.iter().map(|details| details.Timestamp));
    }
    temp_vec.sort();
    temp_vec.dedup();
    temp_vec
}

fn create_analysis_map(vec_Timestamp_sorted: &[i64], performances: &[StockAnalysis]) -> HashMap<i64, Vec<StockAnalysis>> {
    let mut analysis_map: HashMap<i64, Vec<StockAnalysis>> = HashMap::new();
    let mut previous_Equity: HashMap<String, f64> = HashMap::new();

    for &Timestamp in vec_Timestamp_sorted {
        let mut filtered_analysis: Vec<StockAnalysis> = performances.iter()
            .filter(|analysis| analysis.Timestamp == Timestamp)
            .cloned()
            .collect();

        for analysis in &mut filtered_analysis {
            if analysis.Equity == 0.0 {
                if let Some(&prev_Equity) = previous_Equity.get(&analysis.symbol) {
                    analysis.Equity = prev_Equity;
                }
            } else {
                previous_Equity.insert(analysis.symbol.clone(), analysis.Equity);
            }
        }

        analysis_map.insert(Timestamp, filtered_analysis);
    }

    analysis_map
}

fn get_top_gainers_losers(api_response: &ApiResponse, timeframe: usize) -> Vec<(String, f64)> {
    let mut stock_returns: Vec<(String, f64)> = Vec::new();

    for (symbol, data) in &api_response.0 {
        let prices: Vec<f64> = data.Trade_details.iter().map(|pd| if pd.Closes != 0.0 { pd.Closes } else { pd.Position_price }).collect();
        let last_index = prices.len() - 1;
        let query_till_index = last_index.saturating_sub(timeframe);

        if query_till_index < prices.len() {
            let initial_price = prices[query_till_index];
            let final_price = prices[last_index];
            let absolute_return = final_price - initial_price;
            let percentage_return = (absolute_return / initial_price) * 100.0;

            stock_returns.push((symbol.clone(), percentage_return));
        }
    }

    stock_returns.sort_by(|a, b| b.1.partial_cmp(&a.1).unwrap_or(std::cmp::Ordering::Equal));

    stock_returns
}

fn get_portfolio_allocation(api_response: &ApiResponse) -> Vec<(String, f64, f64)> {
    let mut stock_Equity_allocation: Vec<(String, f64, f64)> = Vec::new();
    let mut stock_Equity: Vec<(String, f64)> = Vec::new();

    for (symbol, data) in &api_response.0 {
        let Equity: Vec<f64> = data.Equity.clone();
        let last_index = Equity.len() - 1;
        let last_Equity = Equity[last_index];
        stock_Equity.push((symbol.clone(), last_Equity));
    }

    let sum: f64 = stock_Equity.iter().map(|(_, value)| value).sum();

    for (symbol, Equity) in &stock_Equity {
        let allocation = (Equity / sum) * 100.0;
        stock_Equity_allocation.push((symbol.clone(), allocation, *Equity));
    }
    
    stock_Equity_allocation
}

fn get_total_Equity(api_response: &ApiResponse) -> Vec<f64> {
    let mut total_Equity = Vec::new();
    let max_num_data_points = api_response
        .0
        .values()
        .map(|data| data.Equity.len())
        .max()
        .unwrap_or(0);

    for i in (0..max_num_data_points).rev() {
        let mut sum = 0.0;

        for data in api_response.0.values() {
            if let Some(Equity) = data.Equity.get(data.Equity.len().saturating_sub(max_num_data_points - i)) {
                sum += Equity;
            }
        }

        total_Equity.push(sum);
    }

    total_Equity.reverse();
    total_Equity
}

fn calculate_total_sharpe_ratio(analysis_map: &HashMap<i64, Vec<StockAnalysis>>) -> HashMap<i64, f64> {
    let mut total_sharpe_ratios = HashMap::new();

    for (&Timestamp, analyses) in analysis_map {
        if Timestamp % 300 == 0 {
            let total_Equity: f64 = analyses.iter().map(|analysis| analysis.Equity).sum();
            let total_sharpe: f64 = analyses.iter()
                .map(|analysis| {
                    let weight_in_portfolio = analysis.Equity / total_Equity;
                    weight_in_portfolio * analysis.sharpe_ratio
                })
                .sum();

            total_sharpe_ratios.insert(Timestamp, total_sharpe);
        }
    }

    total_sharpe_ratios
}

fn calculate_total_sortino_ratio(analysis_map: &HashMap<i64, Vec<StockAnalysis>>) -> HashMap<i64, f64> {
    let mut total_sortino_ratios = HashMap::new();

    for (&Timestamp, analyses) in analysis_map {
        if Timestamp % 300 == 0 {
            let total_Equity: f64 = analyses.iter().map(|analysis| analysis.Equity).sum();
            let total_sortino: f64 = analyses.iter()
                .map(|analysis| {
                    let weight_in_portfolio = analysis.Equity / total_Equity;
                    weight_in_portfolio * analysis.sortino_ratio
                })
                .sum();

            total_sortino_ratios.insert(Timestamp, total_sortino);
        }
    }

    total_sortino_ratios
}

async fn connect_to_db() -> Result<Pool<MySql>, Error> {
    let database_url = "mysql://root:Karma100%@tcp(host.docker.internal:3306)/alert_trade_db";
    
    let pool = MySqlPoolOptions::new()
        .max_connections(5)
        .connect(database_url)
        .await?;

    info!("Connected to the database");

    Ok(pool)
}

async fn execute_stored_procedure(pool: &Pool<MySql>, entity_id: i64) -> Result<Vec<String>, Error> {
    info!("Executing stored procedure for entity_id: {}", entity_id);

    let results = sqlx::query(
        "CALL alert_trade_db.stp_Get_Distinct_Symbols_By_User(?)"
    )
    .bind(entity_id)
    .fetch_all(pool)
    .await?;
    
    let mut symbols = Vec::new();

    for row in results {
        let s_ticker: String = row.try_get(0)?;  // Ensure Row trait is imported
        symbols.push(s_ticker);
    }

    info!("Fetched {} symbols for entity_id: {}", symbols.len(), entity_id);

    Ok(symbols)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize the logger with a default level of Info
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
    
    info!("Starting server...");

    HttpServer::new(|| {
        App::new()
            .service(user_analytics_info)
    })
    .bind(("127.0.0.1", 8000))?
    .run()
    .await
}
