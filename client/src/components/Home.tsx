import React, { useEffect, useState } from "react";
import { Layout } from "./Layout";
import type { TArticle } from "../types/articles";
import { ArrowRight, RefreshCw, RotateCcw } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { article } from "../API/articles";
import { useSearchParams } from "react-router";
import { ArticleCard } from "./ArticleCard";
import { InputField } from "./InputField";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Button } from "./Button";
import { SquareBlocks } from "./SquareBlocks";
import type { Loader } from "../types/loader";
import type { Pagination } from "../types/pagination";
import { SearchArticles } from "./SearchArticles";
import { DatePicker } from "./DatePicker";

export const Home: React.FC = () => {
  const [articles, setArticles] = useState<TArticle["article"][]>([]);
  const [pagination, setPagination] = useState<Pagination>();
  const [searchParams, setSearchParams] = useSearchParams();
  const [loader, setLoader] = useState<Loader>("INITIAL");
  const articleIDCursor = searchParams.get("aIDCursor");
  const dateCursor = searchParams.get("dCursor");
  const offset = searchParams.get("offset");

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
        setSearchParams({
          dCursor: new Date(values.timeTravelBitcoin).toISOString(),
          aIDCursor: "",
          query: "",
        });
      } catch (error) {
        helpers.setStatus({ success: false });
        helpers.setSubmitting(false);
      }
    },
  });

  // const articles: TArticle["article"][] = data?.data ?? [];
  // const pagination: Pagination = data?.pagination ?? {};
  const isInitialLoader = loader === "INITIAL";
  const isArticleIDCursorLoader = loader === "ARTICLE_ID_CURSOR";
  const isDateCursorLoader = loader === "DATE_CURSOR";

  const triggerLoadMoreArticles = () => {
    setLoader(() => "ARTICLE_ID_CURSOR");
    setSearchParams({
      dCursor: "",
      aIDCursor: pagination!.prevCursor,
      query: "",
    });
  };

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
              {/* <InputField
                name={"timeTravelBitcoin"}
                type={"date"}
                formik={formik}
              /> */}
              <DatePicker
                name={"timeTravelBitcoin"}
                // type={"date"}
                formik={formik}
              />
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
            onSuccess={function (result: any): void {
              throw new Error("Function not implemented.");
            }}
            onQueryValue={function (hasQueryValue: boolean): void {
              throw new Error("Function not implemented.");
            }}
          />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {articles.map((article, index) => (
            <div key={index} className="w-full">
              <ArticleCard article={article} />
            </div>
          ))}
        </div>
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
        <div>
          <SquareBlocks />
        </div>
      </div>
    </Layout>
  );
};
