import React, { useState, Fragment, type ReactNode, useEffect } from "react";
import ReactDOM from "react-dom";
import { twMerge } from "tailwind-merge";
import { X } from "lucide-react";

interface ModalOverlayProps {
  onClose: () => void;
}

interface ModalContentProps {
  onClose: () => void;
  content: ReactNode;
  className?: string;
}

interface ModalProps {
  openModalElement: ReactNode;
  children: ReactNode;
  className?: string;
  closed?: boolean;
}

const ModalOverlay: React.FC<ModalOverlayProps> = (props) => {
  return (
    <div
      className={`fixed top-0 left-0 z-[70] h-full w-full bg-[rgba(0,0,0,0.5)]
       transition-all`}
      onClick={props.onClose}
    />
  );
};

const ModalContent: React.FC<ModalContentProps> = (props) => {
  return (
    <div
      className={twMerge(
        `relative z-[1000] rounded-md bg-color-bg-primary p-0
          transition-all`,
        props.className
      )}
    >
      <span
        className="absolute right-4 top-1 z-[2000] cursor-pointer"
        onClick={props.onClose}
      >
        <X className="text-[#868e96] w-4 h-4" />
      </span>
      {props.content}
    </div>
  );
};

export const Modal: React.FC<ModalProps> = (props) => {
  const [isOpen, setIsOpen] = useState(false);

  const onOpenHandler = () => setIsOpen(() => !isOpen);
  const onCloseHandler = () => setIsOpen(() => !isOpen);

  const createAppendPortalElement = () => {
    const portalElement = document.createElement("div");
    portalElement.setAttribute("id", "portal");
    const body = document.body;
    body.appendChild(portalElement);
  };
  createAppendPortalElement();

  useEffect(() => {
    const autoCloseHandler = () => {
      const isClosed = props.closed !== undefined && props.closed;
      if (!isClosed) return;
      setIsOpen(() => !isClosed);
    };

    autoCloseHandler();
  }, [props.closed]);

  return (
    <Fragment>
      <>
        <div onClick={() => onOpenHandler()}>{props?.openModalElement}</div>

        {isOpen &&
          ReactDOM.createPortal(
            <div
              className="fixed top-0 left-0 z-[60] flex h-[100vh]
               w-[100vw] items-center justify-center transition-all"
            >
              <ModalOverlay onClose={() => onCloseHandler()} />
              <ModalContent
                content={props.children}
                onClose={() => onCloseHandler()}
                className={props?.className}
              />
            </div>,
            document.getElementById("portal")!
          )}
      </>
    </Fragment>
  );
};
