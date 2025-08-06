import React, { useEffect, useRef } from "react";

interface HeaderBgProps {
  children: React.ReactNode;
}

interface Wave {
  color: string;
  baseAmplitude: number;
  amplitudeOffset: number;
  wavelength: number;
  speed: number;
  phase: number;
  oscillationSpeed: number;
  time: number;
}

export const HeaderBg: React.FC<HeaderBgProps> = ({ children }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let width: number, height: number;
    // const colors = ["#f7931a", "#ff9500", "#00d4aa", "#22c55e"];
    const colors = ["#f7931a", "#f76707", "#00d4aa", "#22c55e"];

    const waves: Wave[] = [];

    const resizeCanvas = () => {
      const dpr = window.devicePixelRatio || 1;
      width = window.innerWidth;
      height = 60;

      canvas.width = width * dpr;
      canvas.height = height * dpr;
      canvas.style.width = width + "px";
      canvas.style.height = height + "px";

      ctx.scale(dpr, dpr);
    };

    resizeCanvas();
    window.addEventListener("resize", resizeCanvas);

    const createWaves = () => {
      colors.forEach((color) => {
        waves.push({
          color,
          baseAmplitude: 10 + Math.random() * 20,
          amplitudeOffset: Math.random() * 5,
          wavelength: (0.005 + Math.random() * 0.01) * 0.75,
          speed: 0.1 + Math.random() * 0.2,
          phase: Math.random() * Math.PI * 2,
          oscillationSpeed: 0.01 + Math.random() * 0.03,
          time: 0,
        });
      });
    };

    let animationId: number;

    const animateWaves = () => {
      ctx.clearRect(0, 0, width, height);

      waves.forEach((wave) => {
        ctx.beginPath();
        ctx.moveTo(0, height / 2);

        const dynamicAmplitude =
          wave.baseAmplitude + Math.sin(wave.time) * wave.amplitudeOffset;

        for (let x = 0; x <= width; x += 10) {
          const y =
            Math.sin(x * wave.wavelength + wave.phase) * dynamicAmplitude +
            height / 2;
          ctx.lineTo(x, y);
        }

        ctx.lineTo(width, height);
        ctx.lineTo(0, height);
        ctx.closePath();

        const gradient = ctx.createLinearGradient(0, 0, width, height);
        gradient.addColorStop(0.85, wave.color);
        gradient.addColorStop(1, "transparent");
        ctx.fillStyle = gradient;
        ctx.fill();

        wave.phase += wave.speed * 0.02;
        wave.time += wave.oscillationSpeed;
      });

      animationId = requestAnimationFrame(animateWaves);
    };

    createWaves();
    animateWaves();

    return () => {
      window.removeEventListener("resize", resizeCanvas);
      if (animationId) {
        cancelAnimationFrame(animationId);
      }
    };
  }, []);

  return (
    <header className="w-full h-16 relative flex items-center justify-center">
      <canvas
        ref={canvasRef}
        className="w-full h-full absolute top-0 left-0 z-[1]"
        style={{ willChange: "transform" }}
      />
      <div
        className="w-full h-full z-10 bg-[rgba(0,0,0,0.25)]
         flex items-center justify-center"
      >
        {children}
      </div>
    </header>
  );
};
