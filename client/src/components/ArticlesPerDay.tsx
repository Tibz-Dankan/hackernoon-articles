import { useQuery } from "@tanstack/react-query";
import React from "react";
import { article } from "../API/articles";
import { RefreshCw } from "lucide-react";
import { ArticleCard } from "./ArticleCard";
import type { TArticle } from "../types/articles";
import { isArrayWithElements } from "../utils/isArrayWithElements";

interface ArticlesPerDayProps {
  day: string;
}

const ArticlesPerDay: React.FC<ArticlesPerDayProps> = (props) => {
  const { isPending, data } = useQuery({
    queryKey: [`articles-${props.day}`],
    queryFn: () => {
      return article.getByDay({
        day: props.day,
      });
    },
  });

  const articles: TArticle["article"][] = data?.data ?? [];
  const hasArticles = isArrayWithElements(articles);

  const getDate = (date: string) => {
    return new Date(date).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  if (isPending) {
    return (
      <div
        className="text-gray-50 flex-1 w-full h-full flex 
        items-center justify-center gap-2 py-4"
      >
        <RefreshCw className="animate-spin" size={24} />
        <span>Loading..</span>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center">
      {hasArticles && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3s gap-4">
          {articles.map((article, index) => (
            <div key={index} className="w-full">
              <ArticleCard article={article} />
            </div>
          ))}
        </div>
      )}
      {!hasArticles && (
        <div className="">
          <p className="text-gray-400 text-sm">{`No Articles for ${getDate(
            props.day
          )}`}</p>
        </div>
      )}
    </div>
  );
};

export default ArticlesPerDay;
