import React from "react";
import type { TArticle } from "../types/articles";
import { formatDate } from "../utils/formatDate";
import { Link } from "react-router";

interface ArticleCardProp {
  article: TArticle["article"];
}

export const ArticleCard: React.FC<ArticleCardProp> = (props) => {
  return (
    <div
      className="w-full p-3 rounded-xl border-[1px] border-[rgba(73,80,87,0.6)]
       space-y-4 bg-[rgba(8,127,91,0.05)] hover:bg-[rgba(8,127,91,0.15)]"
    >
      <div className="w-full inline-block relative">
        <Link
          to={props.article.href!}
          target="_blank"
          rel="noopener noreferrer"
        >
          <img
            src={props.article.imageUrl}
            alt={props.article.title}
            className="w-full h-52 object-cover object-center rounded-lg mx-auto
            bg-(--clr-background)"
          />
        </Link>
        <div
          className="absolute bottom-5 right-5 text-gray-300 px-2
           bg-[rgba(8,127,91,0.25)] text-xl font-bold rounded-3xl
           z-10"
        >
          {props.article.tag}
        </div>
      </div>
      <div className="h-16 text-base text-gray-50 hover:underline font-semibold">
        <Link
          to={props.article.href!}
          target="_blank"
          rel="noopener noreferrer"
        >
          <p>{props.article.title}</p>
        </Link>
      </div>
      <div className="flex justify-start items-start">
        <div className="w-10 h-10 rounded-full border-[1px] border-[rgba(73,80,87,0.45)]">
          <img
            src={props.article.author.avatarUrl}
            alt={props.article.author.avatarFilename}
            className="w-10 h-10 rounded-full"
          />
        </div>
        <div className="ml-2 flex flex-col justify-start">
          <span className="text-gray-50 text-[12px] font-semibold">
            {props.article.author.name}
          </span>
          <span className="text-gray-100 text-[8px]">
            {formatDate(props.article.postedAt)}
          </span>
          <span className="text-gray-100 text-[8px]">
            {props.article.readDuration}
          </span>
        </div>
      </div>
    </div>
  );
};
