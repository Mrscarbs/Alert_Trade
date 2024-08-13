use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::OpenOptions;
use log::info;

#[derive(Serialize, Clone)]
struct RequestBody {
    list_symbols: Vec<String>,
    entity_id: i64,
}

#[derive(Debug, Deserialize, Clone)]
struct PriceDetails {
    closes: f64,
    high: f64,
    low: f64,
    timestamp: i64,
    position_price: f64,
    position_quantity: i64,
}

#[derive(Debug, Deserialize, Clone)]
struct CompanyData {
    equity: Vec<f64>,
    trade_details: Vec<PriceDetails>,
}

#[derive(Debug, Deserialize, Clone)]
struct ApiResponse(HashMap<String, CompanyData>);

#[derive(Debug, Clone)]
struct StockPerformance {
    symbol: String,
    sharpe_ratios: Vec<f64>,
}

#[derive(Debug, Clone)]
struct StockAnalysis {
    symbol: String,
    sharpe_ratio: f64,
    sortino_ratio: f64,
    timestamp: i64,
    equity: f64,
}

// Function to calculate returns from prices
fn calculate_returns(prices: &[f64]) -> Vec<f64> {
    prices.windows(2)
        .map(|w| (w[1] - w[0]) / w[0])
        .collect()
}

// Function to calculate Sharpe Ratio
fn calculate_sharpe_ratio(returns: &[f64]) -> f64 {
    let mean: f64 = returns.iter().sum::<f64>() / returns.len() as f64;
    let std_dev: f64 = (returns.iter().map(|r| (r - mean).powi(2)).sum::<f64>() / returns.len() as f64).sqrt();
    if std_dev == 0.0 {
        return 0.0;
    }
    mean / std_dev
}

// Function to calculate Sortino Ratio
fn calculate_sortino_ratio(returns: &[f64]) -> f64 {
    let mean: f64 = returns.iter().sum::<f64>() / returns.len() as f64;
    let negative_std_dev: f64 = (returns.iter().filter(|r| **r < 0.0).map(|r| r.powi(2)).sum::<f64>() / returns.len() as f64).sqrt();
    if negative_std_dev == 0.0 {
        return 0.0;
    }
    mean / negative_std_dev
}

// Function to calculate cumulative Sharpe and Sortino Ratios for all previous data points starting from the second index
fn cumulative_analysis(returns: &[f64], timestamps: &[i64], equity: &[f64], symbol: &str) -> Vec<StockAnalysis> {
    (2..=returns.len())
        .map(|i| {
            let sharpe_ratio = calculate_sharpe_ratio(&returns[..i]);
            let sortino_ratio = calculate_sortino_ratio(&returns[..i]);
            StockAnalysis {
                symbol: symbol.to_string(),
                sharpe_ratio,
                sortino_ratio,
                timestamp: timestamps[i],
                equity: equity[i],
            }
        })
        .collect()
}

fn all_timestamp_list_creator(stock_map: &ApiResponse) -> Vec<i64> {
    let mut temp_vec = Vec::new();
    for (_key, value) in &stock_map.0 {
        temp_vec.extend(value.trade_details.iter().map(|details| details.timestamp));
    }
    println!("vec: {:?}", temp_vec);
    println!("vec_length: {:?}", temp_vec.len());
    temp_vec.sort();
    temp_vec.dedup(); // Ensure unique timestamps
    temp_vec
}

// Function to create a HashMap of StockAnalysis by timestamp
fn create_analysis_map(vec_timestamp_sorted: &[i64], performances: &[StockAnalysis]) -> HashMap<i64, Vec<StockAnalysis>> {
    let mut analysis_map: HashMap<i64, Vec<StockAnalysis>> = HashMap::new();
    let mut previous_equity: HashMap<String, f64> = HashMap::new();

    for &timestamp in vec_timestamp_sorted {
        let mut filtered_analysis: Vec<StockAnalysis> = performances.iter()
            .filter(|analysis| analysis.timestamp == timestamp)
            .cloned()
            .collect();

        for analysis in &mut filtered_analysis {
            if analysis.equity == 0.0 {
                if let Some(&prev_equity) = previous_equity.get(&analysis.symbol) {
                    analysis.equity = prev_equity;
                }
            } else {
                previous_equity.insert(analysis.symbol.clone(), analysis.equity);
            }
        }

        analysis_map.insert(timestamp, filtered_analysis);
    }

    analysis_map
}

// Function to get top gainers and losers
fn get_top_gainers_losers(api_response: &ApiResponse, timeframe: usize) -> Vec<(String, f64)> {
    let mut stock_returns: Vec<(String, f64)> = Vec::new();

    for (symbol, data) in &api_response.0 {
        let prices: Vec<f64> = data.trade_details.iter().map(|pd| if pd.closes != 0.0 { pd.closes } else { pd.position_price }).collect();
        let last_index = prices.len() - 1; // Last index
        let query_till_index = last_index.saturating_sub(timeframe); // Start of the timeframe

        // Calculate the absolute and percentage return
        if query_till_index < prices.len() {
            let initial_price = prices[query_till_index];
            let final_price = prices[last_index];
            let absolute_return = final_price - initial_price;
            let percentage_return = (absolute_return / initial_price) * 100.0;

            // Store the stock symbol and its percentage return
            stock_returns.push((symbol.clone(), percentage_return));
        }
    }

    // Sort the stocks by percentage return from highest to lowest
    stock_returns.sort_by(|a, b| b.1.partial_cmp(&a.1).unwrap_or(std::cmp::Ordering::Equal));

    stock_returns
}

// Function to get portfolio allocation
fn get_portfolio_allocation(api_response: &ApiResponse) {
    let mut stock_equity_allocation: Vec<(String, f64, &f64)> = Vec::new();
    let mut stock_equity: Vec<(String, f64)> = Vec::new();

    for (symbol, data) in &api_response.0 {
        let equity: Vec<f64> = data.equity.clone();
        let last_index = equity.len() - 1;
        let last_equity = equity[last_index];
        stock_equity.push((symbol.clone(), last_equity));
    }

    let sum: f64 = stock_equity.iter().map(|(_, value)| value).sum();

    for (symbol, equity) in &stock_equity {
        let allocation = (equity / sum) * 100.0;
        println!("Symbol: {}, Allocation: {:.2}%", symbol, allocation);
        stock_equity_allocation.push((symbol.clone(), allocation, equity))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Set up logging to a file
    let log_file = OpenOptions::new()
        .create(true)
        .write(true)
        .truncate(true)
        .open("log.txt")?;
    
    env_logger::Builder::new()
        .filter_level(log::LevelFilter::Info)
        .target(env_logger::Target::Pipe(Box::new(log_file)))
        .init();

    let client = Client::new();

    let body = RequestBody {
        list_symbols: vec!["RELIANCE".to_string(), "TCS".to_string()],
        entity_id: 1,
    };

    let response = client
        .post("http://localhost:8080/get_historical_performance")
        .json(&body)
        .send()
        .await?;

    let response_text = response.text().await?;
    info!("Raw Response: {}", response_text);
    println!("Raw Response: {}", response_text);

    // Trim whitespace to avoid issues with trailing characters
    let trimmed_response_text = response_text.trim();

    // Parse the response as an ApiResponse
    let api_response: ApiResponse = serde_json::from_str(&trimmed_response_text)
        .map_err(|e| Box::new(e) as Box<dyn std::error::Error>)?;

    let mut performances = Vec::new();

    for (symbol, data) in &api_response.0 {
        println!("Symbol: {}", symbol);
        println!("Equity Values: {:?}", data.equity);
        println!("Trade Details: {:?}", data.trade_details);

        // Extract prices from trade details, using position_price if closes is zero
        let prices: Vec<f64> = data.trade_details.iter().map(|pd| if pd.closes != 0.0 { pd.closes } else { pd.position_price }).collect();
        let timestamps: Vec<i64> = data.trade_details.iter().map(|pd| pd.timestamp).collect();
        let equity: Vec<f64> = data.equity.clone();
        
        println!("Prices: {:?}", prices);

        // Calculate returns based on prices
        let returns = calculate_returns(&prices);
        println!("Returns: {:?}", returns);

        // Calculate cumulative Sharpe and Sortino Ratios based on returns
        let analysis = cumulative_analysis(&returns, &timestamps, &equity, symbol);
        performances.extend(analysis);
    }

    // Print the analysis in a readable format
    for performance in &performances {
        println!("Symbol: {}", performance.symbol);
        println!("Timestamp: {}", performance.timestamp);
        println!("Sharpe Ratio: {}", performance.sharpe_ratio);
        println!("Sortino Ratio: {}", performance.sortino_ratio);
        println!("Equity: {}", performance.equity);
    }

    // Create a list of all timestamps
    let vec_timestamp_sorted = all_timestamp_list_creator(&api_response);
    println!("Sorted Timestamps: {:?}", vec_timestamp_sorted);

    // Create the HashMap using the separate function
    let analysis_map: HashMap<i64, Vec<StockAnalysis>> = create_analysis_map(&vec_timestamp_sorted, &performances);

    // Print the HashMap to verify the contents
    for (timestamp, analyses) in &analysis_map {
        println!("Timestamp: {}", timestamp);
        for analysis in analyses {
            println!("  Symbol: {}, Sharpe Ratio: {}, Sortino Ratio: {}, Equity: {}", 
                analysis.symbol, analysis.sharpe_ratio, analysis.sortino_ratio, analysis.equity);
        }
    }

    // Get and print top gainers and losers
    let top_gainers_losers = get_top_gainers_losers(&api_response, 5);
    println!("Top Gainers and Losers:");
    for (symbol, percentage_return) in top_gainers_losers {
        println!("Symbol: {}, Percentage Return: {:.2}%", symbol, percentage_return);
    }

    // Get and print portfolio allocation
    get_portfolio_allocation(&api_response);

    Ok(())
}
