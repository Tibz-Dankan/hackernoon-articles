import type { ReactNode } from "react";
import React, { Fragment } from "react";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 2000,
    },
    mutations: {
      networkMode: "always",
    },
  },
});

interface ReactQueryProviderProps {
  children: ReactNode;
}

export const ReactQueryProvider: React.FC<ReactQueryProviderProps> = (
  props
) => {
  return (
    <Fragment>
      <QueryClientProvider client={queryClient}>
        {props.children}
      </QueryClientProvider>
    </Fragment>
  );
};
