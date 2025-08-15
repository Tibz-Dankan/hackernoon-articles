import React from "react";
import { useFormik } from "formik";
import * as Yup from "yup";
import { useMutation } from "@tanstack/react-query";
import { Search, Loader2 } from "lucide-react";
import { Button } from "./Button";
import type { TArticle } from "../types/articles";
import { article } from "../API/articles";
import { useNavigate, useSearchParams } from "react-router";

interface SearchArticlesProps {
  onSuccess: (result: any) => void;
  onQueryValue: (hasQueryValue: boolean) => void;
}

export const SearchArticles: React.FC<SearchArticlesProps> = (props) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();

  const { isPending, mutate } = useMutation({
    mutationFn: article.search,
    onSuccess: async (response: any) => {
      console.log("response:", response);
      props.onSuccess(response);
    },
    onError: (error: any) => {
      console.log("Error: ", error.message);
    },
  });

  const initialValues: TArticle["searchArticles"] = {
    query: "",
    limit: 18,
  };

  const formik = useFormik({
    initialValues: initialValues,
    validationSchema: Yup.object({
      query: Yup.string().max(255).required("search input is required"),
    }),

    onSubmit: async (values, helpers) => {
      values.query = values.query.trim();

      try {
        mutate(values);
        setSearchParams(
          (prev) => {
            prev.set("query", values.query);
            return prev;
          },
          { replace: true }
        );
      } catch (err: any) {
        helpers.setStatus({ success: false });
        helpers.setSubmitting(false);
        console.log("Error: ", err.message);
      }
    },
  });

  const updateSearchParamsOnChangeHandler = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const query = event.target.value;
    props.onQueryValue(!!event.target.value);

    console.log("query input ", query);

    if (!query) {
      console.log("No query string");
      const currentParams = new URLSearchParams(searchParams.toString());
      currentParams.delete("query");
      navigate(`${location.pathname}?${currentParams.toString()}`);
      return;
    }

    setSearchParams(
      (prev) => {
        prev.set("query", event.target.value);
        return prev;
      },
      { replace: true }
    );
  };

  const hasSearchQuery = !!formik.values["query"];

  return (
    <div className="w-full flex items-center justify-center">
      <div className="relative w-full">
        <form onSubmit={formik.handleSubmit} className="w-full relative">
          <input
            type="text"
            id="query"
            required
            onBlur={formik.handleBlur}
            onChange={(event) => {
              formik.handleChange(event),
                updateSearchParamsOnChangeHandler(event);
            }}
            value={formik.values["query"]}
            placeholder={"Search articles by tag index e.g a23"}
            className="pl-3 py-[10px] pr-9 outline-none rounded-md border-[1px]
            border-[rgba(73,80,87,0.6)] focus:border-[1px] focus:border-(--clr-primary)
            transition-all text-sm w-full focus:outline-none
            focus:shadow-[0px_0px_0px_4px_rgba(12,166,120,0.3)]
            text-gray-50 bg-(--clr-background)s appearance-none"
          />
          <Button
            label={
              <>
                {!isPending && (
                  <Search
                    size={20}
                    className={`${
                      hasSearchQuery ? "text-gray-200" : "text-[#868e96]"
                    }`}
                  />
                )}
                {isPending && (
                  <Loader2
                    size={20}
                    className="text-color-text-primary -mr-2 animate-spin"
                  />
                )}
              </>
            }
            type="submit"
            disabled={isPending}
            className={`px-2 py-[6px] h-auto absolute top-[5px] right-1
              flex items-center justify-center
             ${hasSearchQuery ? "bg-(--clr-primary)" : ""}`}
          />
        </form>
      </div>
    </div>
  );
};
