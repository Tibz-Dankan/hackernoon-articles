import React, { useEffect, useState } from "react";
import { Layout } from "./Layout";
import type { TArticle } from "../types/articles";
import {
  ArrowLeft,
  ArrowRight,
  InfoIcon,
  RefreshCw,
  RotateCcw,
} from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { article } from "../API/articles";
import { useNavigate, useSearchParams } from "react-router";
import { ArticleCard } from "./ArticleCard";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Button } from "./Button";
import type { Loader } from "../types/loader";
import type { Pagination } from "../types/pagination";
import { SearchArticles } from "./SearchArticles";
import { DatePicker } from "./DatePicker";
import { X } from "lucide-react";
import { SquareBlockLayout } from "./SquareBlockLayout";

export const Home: React.FC = () => {
  const [articles, setArticles] = useState<TArticle["article"][]>([]);
  const [searchArticles, setSearchArticles] = useState<TArticle["article"][]>(
    []
  );
  const [pagination, setPagination] = useState<Pagination>();
  const [showArticle, setShowArticle] = useState<boolean>(true);
  const [searchResultCount, setSearchResultCount] = useState<number>();
  // const [disableNextHandler, setDisableNextHandler] = useState(false);
  const [disablePrevHandler, setDisablePrevHandler] = useState(false);

  const navigate = useNavigate();

  const [searchParams, setSearchParams] = useSearchParams();
  const [loader, setLoader] = useState<Loader>("INITIAL");
  const articleIDCursor = searchParams.get("aIDCursor");
  const dateCursor = searchParams.get("dCursor");
  const offset = searchParams.get("offset");
  const searchQuery = searchParams.get("query");

  const { isPending, data } = useQuery({
    queryKey: [`articles-${articleIDCursor}-${dateCursor}-${offset}`],
    queryFn: () => {
      return article.getAll({
        limit: 18,
        articleIDCursor: !!articleIDCursor ? articleIDCursor : "",
        dateCursor: !!dateCursor ? dateCursor : "",
        offset: "",
      });
    },
  });

  const initialValues = {
    timeTravelBitcoin: "",
  };

  const formik = useFormik({
    initialValues: initialValues,
    validationSchema: Yup.object({
      timeTravelBitcoin: Yup.string().max(255).required("Required"),
    }),
    onSubmit: async (values, helpers) => {
      try {
        setLoader(() => "DATE_CURSOR");
        setSearchParams(
          (prev) => {
            prev.set(
              "dCursor",
              new Date(values.timeTravelBitcoin).toISOString()
            );
            prev.set("aIDCursor", "");
            prev.set("query", "");
            return prev;
          },
          { replace: true }
        );
      } catch (error) {
        helpers.setStatus({ success: false });
        helpers.setSubmitting(false);
      }
    },
  });

  const isInitialLoader = loader === "INITIAL";
  const isArticleIDCursorLoader = loader === "ARTICLE_ID_CURSOR";
  const isDateCursorLoader = loader === "DATE_CURSOR";

  const triggerLoadMoreArticles = () => {
    setLoader(() => "ARTICLE_ID_CURSOR");
    setSearchParams(
      (prev) => {
        prev.set("dCursor", "");
        prev.set("aIDCursor", pagination!.prevCursor);
        prev.set("query", "");
        return prev;
      },
      { replace: false }
    );
  };

  const prevLoadHandler = () => {
    if (!articleIDCursor) {
      setDisablePrevHandler(() => true);
      return;
    }
    navigate(-1);
  };

  const closeSearchResultHandler = () => {
    const currentParams = new URLSearchParams(searchParams.toString());
    currentParams.delete("query");
    navigate(`${location.pathname}?${currentParams.toString()}`);

    setSearchResultCount(() => 0);
    setShowArticle(() => true);
  };

  const getArticleHeader = () => {
    if (dateCursor) {
      const date = new Date(dateCursor).toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
      return `Showing articles from ${date} onwards (going back in time)`;
    }

    return "Showing latest articles";
  };

  const onSearchSuccess = (results: any) => {
    console.log("search results:", results);
    setShowArticle(() => false);
    setSearchArticles(() => results.data);
    setSearchResultCount(() => results.data.length);
  };

  useEffect(() => {
    const updateNavButtonDisabilityHandler = () => {
      if (!articleIDCursor) {
        setDisablePrevHandler(() => true);
        return;
      }
      setDisablePrevHandler(() => false);
    };
    updateNavButtonDisabilityHandler();
  }, [articleIDCursor]);

  useEffect(() => {
    const updateArticles = () => {
      setArticles(() => data?.data ?? []);
      setPagination(() => data?.pagination ?? {});
    };
    updateArticles();
  }, [data]);

  if (isPending && isInitialLoader) {
    return (
      <Layout>
        <div
          className="text-gray-50 flex-1 w-full h-full flex 
           items-center justify-center gap-2"
        >
          <RefreshCw className="animate-spin" size={24} />
          <span>Loading..</span>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="w-full flex flex-col items-center gap-8">
        <div className="w-full flex flex-col items-center gap-8 px-4 sm:px-16 lg:px-20">
          <form
            onSubmit={formik.handleSubmit}
            className="w-full flex flex-col items-center gap-2 text-gray-200 self-center"
          >
            <div
              className="w-full flex items-center justify-center gap-2 text-base
              text-gray-400"
            >
              <span className="hidden sm:block">
                <RotateCcw size={20} />
              </span>
              <span className="text-center">
                Time Travel to Our Bitcoin Articles of
              </span>
            </div>
            <div className="w-full flex flex-col sm:flex-row items-center justify-center gap-2">
              <div className="w-full sm:w-auto">
                <DatePicker name={"timeTravelBitcoin"} formik={formik} />
              </div>
              <Button
                label={
                  <>
                    {!isPending && (
                      <div className="flex items-center justify-center gap-2">
                        <span>Submit</span>
                      </div>
                    )}
                    {isPending && isDateCursorLoader && (
                      <div className="flex items-center justify-center gap-2">
                        <RefreshCw className="animate-spin" size={24} />
                        <span>Submitting</span>
                      </div>
                    )}
                  </>
                }
                type={"submit"}
                disabled={isPending}
                className="min-w-40 w-full sm:w-auto"
              />
            </div>
          </form>
          <div className="w-full">
            <SearchArticles
              onSuccess={onSearchSuccess}
              onQueryValue={function (hasQueryValue: boolean): void {
                console.log("hasQueryValue: ", hasQueryValue);
              }}
            />
          </div>
          {showArticle && (
            <div
              className="flex items-center justify-start gap-2 px-4 py-3
               rounded-md bg-[#0ca678]/[0.3] w-full"
            >
              <span className="w-5 h-5">
                <InfoIcon className="w-5 h-5 text-gray-200" />
              </span>
              <span className="text-gray-200">{getArticleHeader()}</span>
            </div>
          )}
          {showArticle && (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {articles.map((article, index) => (
                <div key={index} className="w-full">
                  <ArticleCard article={article} />
                </div>
              ))}
            </div>
          )}
          {showArticle && (
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
                        <span>Next</span>
                        <ArrowRight size={20} />
                      </div>
                    )}
                    {isPending && isArticleIDCursorLoader && (
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
          )}
          {!!searchResultCount && (
            <div
              className="w-full  relative flex items-center justify-between
               gap-4 px-4 py-3 rounded-md bg-[#0ca678]/[0.3]"
            >
              <p className="text-base w-/5 text-gray-200">
                {searchResultCount} search results for "{searchQuery}"
              </p>
              <span
                className="cursor-pointer"
                onClick={() => closeSearchResultHandler()}
              >
                <X className="w-5 h-5 text-gray-200" />
              </span>
            </div>
          )}
          {!!searchResultCount && (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {searchArticles.map((article, index) => (
                <div key={index} className="w-full">
                  <ArticleCard article={article} />
                </div>
              ))}
            </div>
          )}
        </div>
        <div className="w-full">
          <SquareBlockLayout />
        </div>
      </div>
    </Layout>
  );
};
