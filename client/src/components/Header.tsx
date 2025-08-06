import React from "react";
import { HeaderBg } from "./HeaderBg";

export const Header: React.FC = () => {
  return (
    <HeaderBg>
      <span className="text-gray-50 text-2xl font-semibold">
        HACKERNOON INDEX
      </span>
    </HeaderBg>
  );
};
