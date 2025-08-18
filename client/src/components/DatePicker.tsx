import React, { Fragment, type ReactNode } from "react";
import { Calendar } from "lucide-react";

interface DatePickerProps extends React.HTMLAttributes<HTMLInputElement> {
  formik?: any;
  name: string;
  placeholder?: string;
  label?: ReactNode;
  required?: boolean;
  min?: string;
  max?: string;
}

export const DatePicker: React.FC<DatePickerProps> = (props) => {
  const formik = props.formik;
  const label = props.label;
  const name = props.name;
  const placeholder = props.placeholder || "YY-MM-DD";
  const isRequired = props.required ? !!props.required : false;

  const hasError = formik?.errors?.[name] && formik?.touched?.[name];

  const formatDateForDisplay = (dateValue: string) => {
    if (!dateValue) return "";
    const date = new Date(dateValue);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const currentValue = formik?.values?.[name] || "";

  return (
    <Fragment>
      <div className="relative py-2 flex flex-col items-start justify-center gap-1 w-full">
        {label && (
          <label
            className={`text-sm first-letter:uppercase font-[400] ${
              hasError ? "text-red-400" : "text-gray-400"
            }`}
            htmlFor={name}
          >
            {label}
          </label>
        )}

        <div className="w-full relative">
          {/* Hidden native date input */}
          <input
            type="date"
            id={name}
            name={name}
            required={isRequired}
            min={props.min}
            max={props.max}
            onBlur={formik?.handleBlur}
            onChange={formik?.handleChange}
            value={currentValue}
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10
            min-w-40"
            style={{
              // Ensure the input covers the entire custom element
              position: "absolute",
              top: 0,
              left: 0,
              minWidth: "160px",
              width: "100%",
              height: "100%",
              opacity: 0,
              cursor: "pointer",
              zIndex: 10,
            }}
          />

          {/* Custom styled display */}
          <div
            className={`
              relative p-2 rounded-md border-[1px] 
              bg-[rgba(8,127,91,0.15)] text-gray-50
              transition-all text-base w-full min-w-40
              flex items-center justify-between
              cursor-pointer min-h-[42px]
              ${hasError ? "border-red-400" : "border-[rgba(73,80,87,0.6)]"}
              hover:border-opacity-80
              focus-within:border-primary
              ${props.className || ""}
            `}
          >
            <span
              className={`${currentValue ? "text-gray-50" : "text-gray-400"}`}
            >
              {currentValue ? formatDateForDisplay(currentValue) : placeholder}
            </span>

            <Calendar
              size={18}
              className={`
                ml-2 flex-shrink-0 transition-colors
                ${hasError ? "text-red-400" : "text-gray-400"}
              `}
            />
          </div>
        </div>

        {/* Error message */}
        {hasError && (
          <span className="text-red-400 text-sm mt-1">
            {formik.errors[name]}
          </span>
        )}
      </div>
    </Fragment>
  );
};
