type Author = {
  id: string;
  name: string;
  avatarUrl: string;
  avatarFilename: string;
  pageUrl: string;
  createdAt: string;
  updatedAt: string;
};

type Article = {
  id: string;
  authorID: string;
  tag: string;
  title: string;
  href?: string;
  imageUrl: string;
  imageFilename: string;
  postedAt: string;
  readDuration: string;
  createdAt: string;
  updatedAt: string;
  author: Prettify<Author>;
};

type GetAllArticles = {
  limit: number;
  articleIDCursor?: string;
  dateCursor?: string;
  offset?: string;
};

type SearchArticles = {
  limit: number;
  query: string;
  articleIDCursor?: string;
  dateCursor?: string;
  offset?: number;
};

type CountArticle = {
  count: number;
  date: string;
};

export type TArticle = {
  article: Prettify<Article>;
  author: Prettify<Author>;
  getAllArticles: Prettify<GetAllArticles>;
  searchArticles: Prettify<SearchArticles>;
  countArticle: Prettify<CountArticle>;
};
