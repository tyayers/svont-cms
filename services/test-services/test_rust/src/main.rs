// import Rocket
#[macro_use] extern crate rocket;

// add our routes and services modules
mod routes;
mod services;

// import our routes
use routes::content::get_posts;
use routes::content::create_post;
use routes::content::get_current_date;
use routes::content::date_plus_month;

// start the web server and mount our get route at "/api". Can replace /api with anything
// or just leave it as "/" as the default location
#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![create_post, get_posts, get_current_date, date_plus_month])
}

