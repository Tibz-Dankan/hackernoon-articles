import React from "react";

export const Footer: React.FC = () => {
  return (
    <footer
      className="w-full flex items-center justify-center  space-x-1
      text-gray-400 py-8 px-4 sm:px-8 lg:px-12 2xl:px-32"
      style={{
        background: `linear-gradient(to bottom, rgba(116,143,252,0.005),
      rgba(116,143,252,0.05),rgba(116,143,252,0.10),rgba(116,143,252,0.015)`,
      }}
    >
      <span>&copy;</span>
      <span>{new Date().getFullYear()}</span>
      <span>Hackernoon Index.</span>
      <span>All rights reserved.</span>
    </footer>
  );
};
