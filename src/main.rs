use clap::{arg, ArgMatches, Command};
use reqwest::{header, Result};
use serde_json::{json, from_str, Value};

fn cli() -> ArgMatches {
    Command::new("jarbas")
        .version("1.0")
        .about("jarbas is a chatgpt cli for personal use")
        .arg(arg!(--key <VALUE>).short('k').long("key").required(true))
        .arg(arg!(--question <VALUE>).short('q').long("question").required(true))
        .get_matches()
}

async fn get_answer(question: &String, key: &String) -> Result<String> {
    let json_body = json!({
        "model": "gpt-3.5-turbo",
        "messages": [
            {
                "role": "user",
                "content": question
            }
        ]
    });

    let client = reqwest::Client::new();
    let res = client
        .post("https://api.openai.com/v1/chat/completions")
        .header(header::AUTHORIZATION, format!("Bearer {}", key))
        .header(header::CONTENT_TYPE, format!("application/json"))
        .body(json_body.to_string())
        .send()
        .await?
        .text()
        .await?;

    let answer: Value = match from_str(&res) {
        Ok(v) => v,
        Err(_) => Value::Null,
    };

    let content = answer["choices"][0]["message"]["content"].to_string();
    // Todo :: Create a parser that would convert this output to a more visual apealing output for
    // terminal

    Ok(content)
}

#[tokio::main]
async fn main() -> Result<()> {
    let matches = cli();

    println!(
        "question: {}",
        matches.get_one::<String>("question").expect("required")
    );
    println!(
        "answer: {}",
        get_answer(
            matches.get_one::<String>("question").expect("required"),
            matches.get_one::<String>("key").expect("required")
            ).await?
    );

    Ok(())
}
