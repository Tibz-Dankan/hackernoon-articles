import React, { Fragment, type ReactNode } from "react";
import { twMerge } from "tailwind-merge";

interface InputFieldProps extends React.HTMLAttributes<HTMLInputElement> {
  formik?: any;
  name: string;
  type: "text" | "password" | "email" | "number" | "date" | "time";
  placeholder?: string;
  label?: ReactNode;
  required?: boolean;
}

export const InputField: React.FC<InputFieldProps> = (props) => {
  const formik = props.formik;
  const label = props.label;
  const name = props.name;
  const placeholder = props.placeholder ? props.placeholder : "";
  const isRequired = props.required ? !!props.required : false;

  const getFieldType = () => {
    return props.type;
  };

  const hasError = formik.errors[`${name}`] && formik.touched[`${name}`];

  return (
    <Fragment>
      <div
        className="relative py-2 flex flex-col items-start
         justify-center gap-1 w-full"
      >
        {label && (
          <label
            className={`text-sm first-letter:uppercase font-[400]
           ${hasError ? "text-danger-400" : "text-muted-highlight-clr"}`}
          >
            {label}
          </label>
        )}
        <div className="w-full relative">
          <input
            type={getFieldType()}
            id={name}
            required={isRequired}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            value={formik.values[`${name}`]}
            placeholder={placeholder}
            className={twMerge(
              `p-2 outline-none rounded-md border-[1px]
              focus:border-primary text-gray-50 bg-[rgba(8,127,91,0.15)]
              transition-all text-base w-full focus:outline-none
               focus:border-(--clr-primary)  bg-(--clr-background)s ${
                 hasError ? "border-red-400" : "border-[rgba(73,80,87,0.6)]"
               }`,
              props.className
            )}
          />
        </div>
      </div>
    </Fragment>
  );
};
