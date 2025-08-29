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
        background: `linear-gradient(to bottom, rgba(8,127,91,0.015),
        rgba(8,127,91,0.15),rgba(8,127,91,0.015))`,
      }}
    >
      <Header />
      <main
        // className="flex items-center justify-center w-full
        //  h-full px-4 sm:px-16 lg:px-20 bg-green-600x mr-4x sm:mr-0x"
        className="flex items-center justify-center w-full h-full"
      >
        {props.children}
      </main>
      <Footer />
    </div>
  );
};
