import React from "react";
import { Layout } from "./Layout";
import type { TArticle } from "../types/articles";
import { useQuery } from "@tanstack/react-query";
import { article } from "../API/articles";
import { useSearchParams } from "react-router";
import { ArticleCard } from "./ArticleCard";

export const Home: React.FC = () => {
  const [searchParams] = useSearchParams();

  // const [searchParams, setSearchParams] = useSearchParams();

  // const updateParams = () => {
  //   // Set new parameters (replaces existing ones)
  //   setSearchParams({
  //     userId: '123',
  //     category: 'electronics'
  //   });
  // };
  const articleIDCursor = searchParams.get("aIDCursor");
  const dateCursor = searchParams.get("dCursor");
  const offset = searchParams.get("offset");

  const { isPending, isError, data, error } = useQuery({
    queryKey: [`articles-${articleIDCursor}-${dateCursor}-${offset}`],
    queryFn: () => {
      return article.getAll({
        limit: 20,
        articleIDCursor: !!articleIDCursor ? articleIDCursor : "",
        dateCursor: !!dateCursor ? dateCursor : "",
        offset: "",
      });
    },
  });

  const articles: TArticle["article"][] = data?.data ?? [];
  if (isPending) {
    return (
      <Layout>
        <div className="text-green-700">Loading</div>
      </Layout>
    );
  }

  // if (isError ) {
  //   return (
  //     <Layout>
  //     <div className="text-green-700">{error.message}</div>
  //   </Layout>
  //   );
  // }

  return (
    <Layout>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {articles.map((article, index) => (
          <div key={index} className="w-full">
            <ArticleCard article={article} />
          </div>
        ))}
      </div>
    </Layout>
  );
};
