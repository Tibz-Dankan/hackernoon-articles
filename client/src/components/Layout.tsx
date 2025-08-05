import React, { type ReactNode } from "react";
import { Header } from "./Header";
import { Footer } from "./Footer";

interface LayoutProps {
  children: ReactNode;
}

export const Layout: React.FC<LayoutProps> = (props) => {
  return (
    <div
      className="w-full h-full min-w-screen min-h-screen flex flex-col
      justify-start items-center gap-8"
      style={{
        // background: `linear-gradient(to bottom, rgba(116,143,252,0.005),
        //   rgba(116,143,252,0.05), rgba(116,143,252,0.10), rgba(116,143,252,0.015))`,
        background: `linear-gradient(to bottom, rgba(8,127,91,0.015),
        rgba(8,127,91,0.15),rgba(8,127,91,0.015))`,
      }}
    >
      <Header />
      <main className="px-4 sm:px-16 lg:px-20">{props.children}</main>
      <Footer />
    </div>
  );
};
