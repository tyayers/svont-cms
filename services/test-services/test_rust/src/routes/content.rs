use rocket::serde::json::{Json};
use rocket::serde::{Serialize, Deserialize};
use google_cloud_storage::client::Client;
use google_cloud_storage::sign::SignedURLOptions;
use google_cloud_storage::sign::SignedURLMethod;
use google_cloud_storage::http::objects::download::Range;
use google_cloud_storage::http::objects::get::GetObjectRequest;
use google_cloud_storage::http::objects::upload::UploadObjectRequest;
// import services module
use crate::services;

#[derive(Debug, Deserialize, Serialize)]
#[serde(crate = "rocket::serde")]
pub struct Content {
    pub id: String,
    pub createTimestamp: String,
    pub content: String,
    pub author: String,
    pub upvotes: i32
}

#[derive(Debug, Deserialize, Serialize)]
#[serde(crate = "rocket::serde")]
pub struct New_Content {
    pub content: String,
    pub author: String
}

// create a struct to hold our Date data
// need serialize/deserialize to convert to/from JSON
#[derive(Debug, Deserialize, Serialize)]
#[serde(crate = "rocket::serde")]
pub struct Date {
    pub day: u32,
    pub month: u32,
    pub year: i32
}

#[get("/posts")]
pub fn get_posts() -> Json<Date> {
    Json(services::content::get_current_date())
}

#[post("/posts", format = "json", data = "<new_post>")]
pub async fn create_post(new_post: Json<New_Content>) -> Json<Content> {

    let real_post: Content = Content {
        id: String::from("1234"),
        createTimestamp: String::from("test_timestamp"),
        content: new_post.content.clone(),
        author: new_post.author.clone(),
        upvotes: 0
    };

    // Create client.
    // let mut client = Client::default().await.unwrap();

    // // Upload the file
    // let uploaded = client.upload_object(&UploadObjectRequest {
    //     bucket: "cms_tg6qp4dq8".to_string(),
    //     name: "post_1234.json".to_string(),
    //     ..Default::default()
    // }, serde_json::to_string(&real_post).unwrap().as_bytes(), "application/octet-stream", None).await;


    Json(real_post)
}

// create get-current-date route under /date and call get_current_date service which will return a Date object
// route returns a Date object converted to JSON
#[get("/posts/get-current-date")]
pub fn get_current_date() -> Json<Date> {
    Json(services::content::get_current_date())
}

// route will accept data in JSON format and expects a date variable in the function parameters
#[post("/date/date-plus-month", format = "json", data = "<date>")]
pub fn date_plus_month(date: Json<Date>) -> Json<Date> {
    Json(services::content::date_plus_month(date))
}

// this 