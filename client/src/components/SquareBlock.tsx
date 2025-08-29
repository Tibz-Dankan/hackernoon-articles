import React from "react";
import type { TArticle } from "../types/articles";
import { ListOrdered } from "lucide-react";

interface SquareBlockProps {
  articleCount: TArticle["countArticle"];
}

export const SquareBlock: React.FC<SquareBlockProps> = (props) => {
  const getDate = (date: string) => {
    return new Date(date).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  // #69db7c rgb(105,219,124)
  // #51cf66 rgb(81,207,102)
  // #40c057 rgb(64,192,87)
  // #37b24d rgb(55,178,77)
  // #2f9e44 rgb(47,158,68)
  // #2b8a3e rgb(43,138,62)

  return (
    <div
      className="w-full h-auto bg-[#0ca678]s aspect-[1/1] p-4 
      flex flex-col justify-between text-gray-50 rounded-mdx roundedx"
      // style={{
      //   background: `linear-gradient(to bottom, rgba(116,143,252,0.005),
      // rgba(116,143,252,0.05),rgba(116,143,252,0.10),rgba(116,143,252,0.015)`,
      // }}

      // style={{
      //   background: `linear-gradient(to bottom, rgb(105,219,124),
      //   rgb(81,207,102),rgb(64,192,87),rgb(55,178,77),rgb(47,158,68),rgb(43,138,62)`,
      // }}

      // style={{
      //   background: `linear-gradient(to bottom, rgba(105,219,124,0.5),
      //   rgba(81,207,102,0.5),rgba(64,192,87,0.5),rgba(55,178,77),rgba(47,158,68,0.5),rgba(43,138,62,0.5)`,
      // }}

      style={{
        background: `linear-gradient(to bottom,
        rgba(105,219,124,0.5),
        rgba(81,207,102,0.5),
        rgba(64,192,87,0.5),
        rgba(55,178,77,0.5),
        rgba(47,158,68,0.5),
        rgba(43,138,62,0.5),
        rgba(43,138,62,0.5),
        rgba(43,138,62,0.5),
        rgba(43,138,62,0.5)
        `,
      }}
    >
      <span>{getDate(props.articleCount.date)}</span>
      <div className="flex items-center justify-between gap-2">
        <ListOrdered size={14} className="text-gray-200" />
        <span className="text-sm text-gray-200">
          {props.articleCount.count} articles
        </span>
      </div>
    </div>
  );
};
