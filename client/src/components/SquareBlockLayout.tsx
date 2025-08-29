import { useQuery } from "@tanstack/react-query";
import React, { useEffect, useState } from "react";
import { article } from "../API/articles";
import { useNavigate, useSearchParams } from "react-router";
import { ArrowLeft, ArrowRight, RefreshCw } from "lucide-react";
import type { TArticle } from "../types/articles";
import { SquareBlock } from "./SquareBlock";
import { Button } from "./Button";
import type { Loader } from "../types/loader";
import { isArrayWithElements } from "../utils/isArrayWithElements";

export const SquareBlockLayout: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [loader, setLoader] = useState<Loader>("INITIAL");
  const dateCursor = searchParams.get("bdCursor");
  // const [disableNextHandler, setDisableNextHandler] = useState(false);
  const [disablePrevHandler, setDisablePrevHandler] = useState(false);
  const navigate = useNavigate();

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
      { replace: false }
    );
  };

  const prevLoadHandler = () => {
    if (!dateCursor) {
      setDisablePrevHandler(() => true);
      return;
    }
    navigate(-1);
  };

  useEffect(() => {
    const updateNavButtonDisabilityHandler = () => {
      if (!dateCursor) {
        setDisablePrevHandler(() => true);
        return;
      }

      setDisablePrevHandler(() => false);
    };
    updateNavButtonDisabilityHandler();
  }, [dateCursor]);

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

  const nextLoadHandler = () => {
    triggerLoadMoreArticles();
  };

  return (
    <div className="w-full space-y-8">
      <div className="">
        {hasArticles && (
          <div className="w-full overflow-x-auto scrollbar-none">
            <div className="flex gap-4 pb-4">
              {articleBlocks?.map((articleCount, index) => (
                <div key={index} className="flex-shrink-0">
                  <SquareBlock articleCount={articleCount} />
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
      <div className="flex items-center justify-center text-(--clr-primary) gap-4">
        <Button
          label={
            <div className="flex items-center justify-center gap-2">
              <ArrowLeft size={20} className="text-inherit" />
              <span className="text-inherit">Previous</span>
            </div>
          }
          type={"button"}
          disabled={disablePrevHandler}
          className="min-w-40 bg-(--clr-background) border-[1px]
           border-[rgba(73,80,87,0.6)] disabled:text-gray-100/50"
          onClick={() => prevLoadHandler()}
        />
        <Button
          label={
            <>
              {!isPending && (
                <div className="flex items-center justify-center gap-2">
                  <span className="text-inherit">Next</span>
                  <ArrowRight size={20} className="text-inherit" />
                </div>
              )}
              {isPending && isDateCursorLoader && (
                <div className="flex items-center justify-center gap-2">
                  <RefreshCw className="animate-spin text-inherit" size={24} />
                  <span className="text-inherit">Loading...</span>
                </div>
              )}
            </>
          }
          type={"button"}
          disabled={isPending}
          className="min-w-40 bg-(--clr-background) border-[1px]
           border-[rgba(73,80,87,0.6)] disabled:text-gray-100/50"
          onClick={() => nextLoadHandler()}
        />
      </div>
    </div>
  );
};
