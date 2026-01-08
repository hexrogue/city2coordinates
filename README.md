# City to Coordinates (Malaysia)

This repository converts Malaysian cities and districts into latitude and longitude coordinates.

It uses the OpenStreetMap Nominatim API to search locations and stores the results in a SQLite database.
This project exists to support other systems that need geographic coordinates, such as prayer time calculation.

Simple goal, clear output, no drama.

## Why This Repo Exists

Instead of mixing data fetching logic inside the preyer time calculator, this repository handles location data only.

Separation like this makes the system easier to maintain, test, and explain in my portfolio.

One repo solves one problem.

## Features

- Converts Malaysian cities and districts into coordinates
- Uses Nominatim OpenStreetMap API
- Saves data into SQLite database
- Rate limited to avoid API abuse
- Structured by state and city

## Tech Stack

- Go
- SQLite
- OpenStreetMap Nominatim API

## Database Schema

Table Name: `zones`

| Column | Type    |
| ------ | ------- |
| id     | INTEGER |
| city   | TEXT    |
| state  | TEXT    |
| lat    | REAL    |
| lon    | REAL    |

## How It Works

1. A predefined list of Malaysian states and cities is used
2. Each city is queried using Nominatim API
3. The first result is selected
4. Latitude and longitude are stored in zones.db
5. A 1 second delay is added between requests to respect API usage policy

## Usage

***1. Install dependencies***

Make sure you have Go installed and SQLite available.

***2. Run The Program***

```bash
go run main.go
```

***3. Output***

A SQLite named `zones.db` will be created with all cities coordinates.

### Example Query

```sql
SELECT city, state, lan, lon FROM zones WHERE state = 'Selangor';
```

## Notes

- This project is intended for data preparation and learning purposes
- Accuracy depends on Nominatim search results
- Not designed for real time geocoding at scale

