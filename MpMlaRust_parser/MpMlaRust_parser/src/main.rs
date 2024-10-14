use csv::{ReaderBuilder, Writer};
use regex::Regex;
use std::error::Error;
use std::fs::File;
use sqlx::{Pool, MySql};
use sqlx::mysql::MySqlPoolOptions;
use tokio;

#[derive(Debug)]
struct Shareholding {
    candidate: String,
    share_name: String,
    amount: i64,
}

fn parse_amount(amount_str: &str) -> i64 {
    amount_str.replace(",", "").parse().unwrap_or(0)
}

async fn connect_to_mysql_db() -> Result<Pool<MySql>, Box<dyn Error>> {
    // Fetching the database URL from the environment variable
    let database_url = std::env::var("DATABASE_URL")?;
    let pool = MySqlPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await?;
    Ok(pool)
}

async fn upsert_shareholding(
    pool: &Pool<MySql>,
    candidate: &str,
    share_name: &str,
    amount: i64
) -> Result<(), Box<dyn Error>> {
    sqlx::query("CALL upsert_shareholding(?, ?, ?)")
        .bind(candidate)
        .bind(share_name)
        .bind(amount)
        .execute(pool)
        .await?;
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    dotenv::dotenv().ok();  // Load environment variables from .env
    let pool = connect_to_mysql_db().await?;

    let input_file = File::open("alert_trade_mp_mla_data.csv")?;
    let mut reader = ReaderBuilder::new()
        .has_headers(true)
        .flexible(true)
        .from_reader(input_file);

    let output_file = File::create("shareholdings_output.csv")?;
    let mut writer = Writer::from_writer(output_file);
    writer.write_record(&["Candidate", "Share Name", "Amount"])?;

    let amount_regex = Regex::new(r"(\d+(?:,\d+)*)")?;
    let share_regex = Regex::new(r"(\d+)\s*shares?\s+of\s+(.+?)(?:\s+having|\s+No\.|\s*$)")?;
    let exclude_regex = Regex::new(r"(Fixed Deposit|A/C No|account number)")?;

    let mut current_candidate = String::new();
    let mut current_share_name = String::new();
    let mut current_amount = String::new();

    for result in reader.records() {
        let record = result?;
        if record.len() < 2 {
            continue;
        }

        let candidate = record[0].trim();
        if !candidate.is_empty() && !candidate.starts_with("Rs") {
            current_candidate = candidate.to_string();
        }

        let description = record[1].trim();
        for line in description.lines() {
            let line = line.trim();

            if exclude_regex.is_match(line) {
                continue;
            }

            if let Some(captures) = share_regex.captures(line) {
                if !current_share_name.is_empty() && !current_amount.is_empty() {
                    let amount_value = parse_amount(&current_amount);
                    writer.write_record(&[&current_candidate, &current_share_name, &current_amount])?;

                    upsert_shareholding(&pool, &current_candidate, &current_share_name, amount_value).await?;
                    current_share_name.clear();
                    current_amount.clear();
                }
                current_share_name = captures[2].trim().to_string();
                current_amount = captures[1].to_string();
            } else if let Some(captures) = amount_regex.captures(line) {
                current_amount = captures[1].to_string();
                if !current_share_name.is_empty() {
                    let amount_value = parse_amount(&current_amount);
                    writer.write_record(&[&current_candidate, &current_share_name, &current_amount])?;

                    upsert_shareholding(&pool, &current_candidate, &current_share_name, amount_value).await?;
                    current_share_name.clear();
                    current_amount.clear();
                }
            } else if !line.is_empty() {
                if !current_share_name.is_empty() {
                    current_share_name += " ";
                }
                current_share_name += line;
            }
        }
    }

    writer.flush()?;
    println!("Conversion completed. Output written to shareholdings_output.csv");

    Ok(())
}
