import { SERVER_URL } from "../constants";
import type { TArticle } from "../types/articles";

class ArticleService {
  getAll = async ({
    limit,
    articleIDCursor,
    dateCursor,
    offset,
  }: TArticle["getAllArticles"]) => {
    const response = await fetch(
      `${SERVER_URL}/api/v0.1/articles?limit=${limit}&articleIDCursor=${articleIDCursor}&dateCursor=${dateCursor}&offset=${offset}`,
      {
        method: "GET",
        headers: {
          "Content-type": "application/json",
        },
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }
    return await response.json();
  };

  search = async ({
    query,
    limit,
    articleIDCursor,
    dateCursor,
    offset,
  }: TArticle["searchArticles"]) => {
    const response = await fetch(
      `${SERVER_URL}/api/v0.1/articles/search?limit=${limit}&query=${query}
      &articleIDCursor=${articleIDCursor}&dateCursor=${dateCursor}&offset=${offset}`,
      {
        method: "GET",
        headers: {
          "Content-type": "application/json",
        },
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    return await response.json();
  };
}

export const article = new ArticleService();
