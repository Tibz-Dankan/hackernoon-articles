import React from "react";
import type { TArticle } from "../types/articles";
import { ListOrdered } from "lucide-react";
import { Modal } from "./Modal";
import ArticlesPerDay from "./ArticlesPerDay";
import { dateToIsoString } from "../utils/dateToIsoString";

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

  return (
    <Modal
      openModalElement={
        <div
          className="w-full h-auto bg-[#0ca678]s aspect-[1/1] p-4 
          flex flex-col justify-between text-gray-50"
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
      }
    >
      <div
        className="bg-(--clr-background) w-[60vw] max-w-[60vw] max-h-[80vh]
         p-4 rounded-md overflow-x-hidden"
      >
        <p
          className="text-gray-50 text-center bg-[#0ca678]/[0.3] p-2 
          rounded mt-2 mb-4"
        >
          {`${getDate(props.articleCount.date)}`}
        </p>
        <ArticlesPerDay day={dateToIsoString(props.articleCount.date)} />
      </div>
    </Modal>
  );
};
