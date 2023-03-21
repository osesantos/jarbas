use std::{fs, io::Write};
use std::io::Read;

use clap::{arg, ArgMatches, Command};
use reqwest::{header, Result};
use serde_json::{from_str, json, Value};

fn cli() -> ArgMatches {
    Command::new("jarbas")
        .version("1.0")
        .about("jarbas is a chatgpt cli for personal use")
        .arg(arg!(--key <VALUE>).short('k').long("key").required(false))
        .arg(
            arg!(--question <VALUE>)
                .short('q')
                .long("question")
                .required(false),
        )
        .arg(arg!(--config).short('c').long("config").required(false))
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
        .await.expect("Error: Unable to perform request")
        .text()
        .await.expect("Error: Unable to obtain result");

    let answer: Value = from_str(&res).expect("Error: Unable to parse result");

    let content = answer["choices"][0]["message"]["content"].to_string();
    // Todo :: Create a parser that would convert this output to a more visual apealing output for
    // terminal

    Ok(content)
}

fn get_config_file() -> std::io::Result<String> {
    let mut home = match dirs::home_dir() {
        Some(dir) => dir,
        None => std::path::PathBuf::new(),
    };

    home.push(".jarbasrc");
    Ok(home.as_os_str().to_str().unwrap().to_string())
}

fn config() -> std::io::Result<()> {
    println!("Welcome to the config...");
    println!("Please provide me an API key to use on every call: ");

    let mut line = String::new();
    std::io::stdin()
        .read_line(&mut line)
        .expect("Failed to read key");

    let file_path = get_config_file()?;
    let borrowed_path = file_path.clone();

    let mut file = fs::File::create(file_path)?;
    file.write(format!("api-key: {}", line).as_bytes())?;

    file.sync_all()?;

    println!("config file \"{}\" created and key added.", borrowed_path);

    Ok(())
}

fn get_key() -> std::io::Result<String> {
    let file_path = get_config_file()?;
    println!("Getting key from {}", file_path);
    let mut file = match fs::File::open(file_path) {
        Ok(file) => file,
        Err(err) => match err.kind() {
            std::io::ErrorKind::NotFound => match fs::File::create(get_config_file()?) {
                Ok(fc) => fc,
                Err(e) => panic!("Problem creating the file: {:?}", e),
            },
            other_error => {
                panic!("Problem opening the file: {:?}", other_error);
            }
        }
    };
    let mut key = String::new();
    file.read_to_string(&mut key)?;

    key = key.trim_start_matches("api-key: ").to_string();
    key = key.trim_end_matches("\n").to_string();

    Ok(key)
}

async fn get_answer_wrapper(matches_key: String, key: String, matches: ArgMatches) -> String {
    let mut answer = String::new();
    if matches_key != "" {
        answer = get_answer(
            matches.get_one::<String>("question").expect("required"),
            &matches_key,
        ).await.expect("Error: Unable to get answer");
    } else if key != "" {
        answer = get_answer(
            matches.get_one::<String>("question").expect("required"),
            &key,
        ).await.expect("Error: Unable to get answer");
    } else {
        answer = "".to_string();
    }
    answer
}

#[tokio::main]
async fn main() -> Result<()> {
    let matches = cli();

    let key = get_key().expect("Unable to get key from file");

    if matches.get_flag("config") {
        match config() {
            Err(e) => panic!("Error: {:?}", e),
            _ => std::process::exit(0)
        }
    }

    let matches_key = String::new();
    matches.get_one::<String>("key").unwrap_or(&matches_key);


    println!(
        "question: {}",
        matches.get_one::<String>("question").expect("required")
    );
    println!(
        "answer: {:#?}",
        get_answer_wrapper(matches_key, key, matches).await
    );

    Ok(())
}
