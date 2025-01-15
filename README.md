# Custom URL Shortener

This is a simple **URL Shortener** service that allows users to shorten long URLs into compact, shareable links. It includes features like rate limiting, custom short URL creation, and expiry times for shortened URLs. This project is built using **Go**, **Redis**, and **Fiber**.

## Features

- **URL Shortening**: Shorten long URLs into short, shareable links.
- **Custom Short URLs**: Users can specify their own custom short URL.
- **Rate Limiting**: Limits the number of requests a user can make per time period.
- **Expiry**: Set an expiry date for shortened URLs.
- **Redis**: Uses Redis to store URL mappings and track rate limits.

## Technologies Used

- **Go**: Backend programming language.
- **Fiber**: Web framework for Go.
- **Redis**: In-memory data store for URL mappings and rate limiting.
- **UUID**: For generating unique short URL identifiers.
- **govalidator**: For validating URLs.

## Installation

### Prerequisites

- Go (1.16+)
- Redis server (locally or using a cloud-based Redis service)

### Step-by-Step Guide

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/your-username/url-shortener.git
   cd url-shortener
