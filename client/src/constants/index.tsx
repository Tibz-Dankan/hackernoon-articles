let SERVER_URL: string;

if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
  SERVER_URL = "http://localhost:3000";
} else {
  SERVER_URL = "https://hackernoon-articles.onrender.com";
}

export { SERVER_URL };
