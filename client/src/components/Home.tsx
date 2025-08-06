import React from "react";
import { Layout } from "./Layout";
import type { TArticle } from "../types/articles";
import { RefreshCw, RotateCcw } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { article } from "../API/articles";
import { useSearchParams } from "react-router";
import { ArticleCard } from "./ArticleCard";
import { InputField } from "./InputField";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Button } from "./Button";
import { SquareBlocks } from "./SquareBlocks";

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
        // submit functionality here
      } catch (error) {
        helpers.setStatus({ success: false });
        helpers.setSubmitting(false);
      }
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
      <div className="flex flex-col gap-8">
        {/* memepool square blocks here */}
        <form
          onSubmit={formik.handleSubmit}
          className="flex flex-col items-center gap-2 text-gray-200 self-center"
        >
          <div
            className="flex items-center justify-center gap-2 text-base
            text-gray-400"
          >
            <span>
              <RotateCcw size={20} />
            </span>
            <span>Time Travel to Our Bitcoin Articles of</span>
          </div>
          <div className="flex items-center justify-center gap-2">
            <InputField
              name={"timeTravelBitcoin"}
              type={"date"}
              formik={formik}
            />
            <Button
              label={
                <>
                  {!isPending && (
                    <div className="flex items-center justify-center gap-2 ">
                      <span>Submit</span>
                    </div>
                  )}
                  {isPending && (
                    <div className="flex items-center justify-center gap-2 ">
                      <RefreshCw className="animate-spin" size={24} />
                      <span>Submitting</span>
                    </div>
                  )}
                </>
              }
              type={"submit"}
              disabled={isPending}
              className="min-w-40"
            />
          </div>
        </form>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {articles.map((article, index) => (
            <div key={index} className="w-full">
              <ArticleCard article={article} />
            </div>
          ))}
        </div>
        <div>
          <SquareBlocks />
        </div>
      </div>
    </Layout>
  );
};
