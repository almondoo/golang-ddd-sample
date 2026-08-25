"use client";

import { motion, useReducedMotion, type HTMLMotionProps } from "motion/react";

type RevealProps = HTMLMotionProps<"div"> & {
  /** whileInViewが再トリガーする余白。CSSのmargin記法(例: "-10%") */
  margin?: `${number}%` | `${number}px`;
};

/**
 * スクロールで画面に入ったらふわっと(fade + rise)現れる汎用ラッパー。
 * OSの「視差効果を減らす」設定が有効なときはopacityのみのフェードにする。
 */
export function Reveal({
  children,
  className,
  margin = "-10%",
  ...props
}: RevealProps) {
  const shouldReduceMotion = useReducedMotion();

  return (
    <motion.div
      className={className}
      initial={shouldReduceMotion ? { opacity: 0 } : { opacity: 0, y: 16 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: margin as never }}
      transition={{ duration: 0.4, ease: "easeOut" }}
      {...props}
    >
      {children}
    </motion.div>
  );
}
