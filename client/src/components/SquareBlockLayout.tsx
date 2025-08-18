import { useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { article } from "../API/articles";
import { useSearchParams } from "react-router";
import { ArrowRight, RefreshCw } from "lucide-react";
import type { TArticle } from "../types/articles";
import { SquareBlock } from "./SquareBlock";
import { Button } from "./Button";
import type { Loader } from "../types/loader";
import { isArrayWithElements } from "../utils/isArrayWithElements";

export const SquareBlockLayout: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [loader, setLoader] = useState<Loader>("INITIAL");
  const dateCursor = searchParams.get("bdCursor");

  const { isPending, data } = useQuery({
    queryKey: [`block-articles-${dateCursor}`],
    queryFn: () => {
      return article.getDayCount({
        limit: 20,
        dateCursor: !!dateCursor ? dateCursor : "",
      });
    },
  });

  const articleBlocks = data?.data as TArticle["countArticle"][];
  const pagination = data?.pagination;
  const hasArticles = isArrayWithElements(articleBlocks);

  const isInitialLoader = loader === "INITIAL";
  const isDateCursorLoader = loader === "B_DATE_CURSOR";

  const triggerLoadMoreArticles = () => {
    setLoader(() => "B_DATE_CURSOR");
    setSearchParams(
      (prev) => {
        prev.set("bdCursor", pagination!.nextCursor);
        return prev;
      },
      { replace: true }
    );
  };

  if (isPending && isInitialLoader) {
    return (
      <div
        className="text-gray-50 flex-1 w-full h-full flex 
        items-center justify-center gap-2"
      >
        <RefreshCw className="animate-spin" size={24} />
        <span>Loading..</span>
      </div>
    );
  }

  return (
    <div className="w-full space-y-8">
      {hasArticles && (
        <div
          className="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-6 
          xl:grid-cols-8 gap-4"
        >
          {articleBlocks?.map((articleCount, index) => (
            <div key={index} className="w-full">
              <SquareBlock articleCount={articleCount} />
            </div>
          ))}
        </div>
      )}
      <div className="flex items-center justify-center text-(--clr-primary)">
        <Button
          label={
            <>
              {!isPending && (
                <div className="flex items-center justify-center gap-2">
                  <span>Next</span>
                  <ArrowRight size={20} />
                </div>
              )}
              {isPending && isDateCursorLoader && (
                <div className="flex items-center justify-center gap-2">
                  <RefreshCw className="animate-spin" size={24} />
                  <span>Loading...</span>
                </div>
              )}
            </>
          }
          type={"button"}
          disabled={isPending}
          className="min-w-40 bg-(--clr-background) border-[1px]
           border-[rgba(73,80,87,0.6)]"
          onClick={() => triggerLoadMoreArticles()}
        />
      </div>
    </div>
  );
};
