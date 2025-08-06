import React from "react";

export const SquareBlocks: React.FC = () => {
  const Block2D = () => (
    <div className="w-16 h-16" style={{ backgroundColor: "#0ca678" }}></div>
  );

  return (
    <div className="flex flex-col items-center gap-4 mt-8">
      {/* First block arrangement: 2|3 */}
      <div className="flex flex-col items-center gap-4">
        <div className="bg-white p-2">
          <div className="flex items-center gap-2">
            {/* Left side: 2 blocks */}
            <div className="flex gap-2">
              <Block2D />
              <Block2D />
            </div>

            {/* Divider */}
            <div className="border-l-2 border-dashed h-20" />

            {/* Right side: 3 blocks */}
            <div className="flex gap-2">
              <Block2D />
              <Block2D />
              <Block2D />
            </div>
          </div>
        </div>
      </div>

      {/* Second block arrangement: 3|2 */}
      <div className="flex flex-col items-center gap-4">
        <div className="bg-white p-2 rounded-lgss shadow-lgs">
          <div className="flex items-center gap-2">
            {/* Left side: 3 blocks */}
            <div className="flex gap-2">
              <Block2D />
              <Block2D />
              <Block2D />
            </div>

            {/* Divider */}
            <div className="border-l-2 border-dashed h-20" />

            {/* Right side: 2 blocks */}
            <div className="flex gap-2">
              <Block2D />
              <Block2D />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
