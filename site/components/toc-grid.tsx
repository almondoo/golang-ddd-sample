"use client";

import Link from "next/link";
import { motion, useReducedMotion, type Variants } from "motion/react";
import type { ReactNode } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { CONCEPT_PAGES, type ConceptPage, pageHref } from "@/lib/pages-data";

/** CONCEPT_PAGESの出現順を保ったまま、連続するgroupごとに束ねる */
function groupConceptPages(pages: ConceptPage[]): { group: ConceptPage["group"]; pages: ConceptPage[] }[] {
  const sections: { group: ConceptPage["group"]; pages: ConceptPage[] }[] = [];
  for (const page of pages) {
    const last = sections[sections.length - 1];
    if (last && last.group === page.group) {
      last.pages.push(page);
    } else {
      sections.push({ group: page.group, pages: [page] });
    }
  }
  return sections;
}

const containerVariants: Variants = {
  hidden: {},
  visible: {
    transition: { staggerChildren: 0.06 },
  },
};

const cardVariants: Variants = {
  hidden: { opacity: 0, y: 16 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: "easeOut" } },
};

const cardVariantsReduced: Variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1, transition: { duration: 0.4, ease: "easeOut" } },
};

/** 目次カードのアイコン。ページごとにトップページのオリジナルSVGを移植する */
const ICONS: Record<string, ReactNode> = {
  domain: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="42" cy="42" r="32" fill="#2dd4bf" opacity="0.4" />
      <circle cx="42" cy="42" r="32" fill="none" stroke="#0d9488" strokeWidth="3" />
      <rect x="24" y="24" width="14" height="14" rx="4" fill="#0d9488" />
      <rect x="46" y="24" width="14" height="14" rx="4" fill="#0d9488" />
      <rect x="24" y="46" width="14" height="14" rx="4" fill="#0d9488" />
      <rect x="46" y="46" width="14" height="14" rx="4" fill="#0d9488" />
    </svg>
  ),
  "context-map": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="42" cy="42" r="14" fill="#0d9488" />
      <circle cx="42" cy="16" r="8" fill="#3E6FB0" />
      <circle cx="68" cy="30" r="8" fill="#5B8F4E" />
      <circle cx="68" cy="54" r="8" fill="#C4577B" />
      <circle cx="16" cy="54" r="8" fill="#B87A1E" />
      <line x1="42" y1="30" x2="42" y2="22" stroke="#0d9488" strokeWidth="3" />
      <line x1="53" y1="35" x2="62" y2="31" stroke="#0d9488" strokeWidth="3" />
      <line x1="53" y1="49" x2="62" y2="53" stroke="#0d9488" strokeWidth="3" />
      <line x1="31" y1="49" x2="22" y2="53" stroke="#0d9488" strokeWidth="3" />
    </svg>
  ),
  "query-service": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="10" y="30" width="26" height="40" rx="8" fill="#3E6FB0" />
      <rect x="48" y="30" width="26" height="40" rx="8" fill="#2dd4bf" />
      <circle cx="61" cy="46" r="6" fill="none" stroke="#fff" strokeWidth="2.6" />
      <line x1="65.2" y1="50.2" x2="70" y2="55" stroke="#fff" strokeWidth="2.6" strokeLinecap="round" />
    </svg>
  ),
  "ubiquitous-language": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="10" y="16" width="64" height="42" rx="16" fill="#E8672E" />
      <polygon points="26,58 40,58 30,72" fill="#E8672E" />
      <circle cx="30" cy="37" r="4.5" fill="#fff" />
      <circle cx="42" cy="37" r="4.5" fill="#fff" />
      <circle cx="54" cy="37" r="4.5" fill="#fff" />
    </svg>
  ),
  "bounded-context": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="8" y="8" width="32" height="32" rx="8" fill="#B87A1E" />
      <rect x="44" y="8" width="32" height="32" rx="8" fill="#3E6FB0" />
      <rect x="8" y="44" width="32" height="32" rx="8" fill="#5B8F4E" />
      <rect x="44" y="44" width="32" height="32" rx="8" fill="#C4577B" />
    </svg>
  ),
  entity: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="42" cy="30" r="18" fill="#C4577B" />
      <circle cx="35" cy="28" r="2.6" fill="#2B2117" />
      <circle cx="49" cy="28" r="2.6" fill="#2B2117" />
      <path
        d="M 34 36 Q 42 42 50 36"
        stroke="#2B2117"
        strokeWidth="2.6"
        fill="none"
        strokeLinecap="round"
      />
      <rect x="18" y="56" width="48" height="20" rx="10" fill="#C4577B" />
      <text x="42" y="70" textAnchor="middle" fill="#fff" fontSize="11" fontWeight="700">
        ID
      </text>
    </svg>
  ),
  "value-object": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="34" cy="42" r="24" fill="#F6E7CC" stroke="#B87A1E" strokeWidth="3" />
      <circle cx="54" cy="42" r="24" fill="#F6E7CC" stroke="#B87A1E" strokeWidth="3" />
      <text x="34" y="48" textAnchor="middle" fill="#B87A1E" fontWeight="800" fontSize="16">
        ¥
      </text>
      <text x="54" y="48" textAnchor="middle" fill="#B87A1E" fontWeight="800" fontSize="16">
        ¥
      </text>
    </svg>
  ),
  aggregate: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="14" y="34" width="56" height="36" rx="10" fill="#FBE4D8" stroke="#E8672E" strokeWidth="3" />
      <polygon points="14,34 30,18 54,18 70,34" fill="#f7d3bb" stroke="#E8672E" strokeWidth="2" />
      <rect x="30" y="42" width="24" height="20" rx="5" fill="#E8672E" />
    </svg>
  ),
  repository: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="16" y="36" width="52" height="34" rx="6" fill="#E9E1F3" stroke="#7A5CA8" strokeWidth="3" />
      <polygon points="16,36 42,16 68,36" fill="#d9c7ee" stroke="#7A5CA8" strokeWidth="2" />
      <rect x="34" y="50" width="16" height="20" fill="#7A5CA8" />
    </svg>
  ),
  "onion-architecture": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="42" cy="42" r="34" fill="#DCE7F5" />
      <circle cx="42" cy="42" r="23" fill="#FBE4D8" />
      <circle cx="42" cy="42" r="12" fill="#E8672E" />
    </svg>
  ),
  "place-order": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="20" cy="42" r="14" fill="#3E6FB0" />
      <text x="20" y="47" textAnchor="middle" fill="#fff" fontWeight="800" fontSize="14">
        1
      </text>
      <polygon points="38,42 58,42 58,34 72,42 58,50 58,42" fill="#E8672E" />
    </svg>
  ),
  "shared-kernel": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="14" y="38" width="56" height="32" rx="7" fill="#0d9488" />
      <rect x="14" y="38" width="56" height="10" rx="5" fill="#0f766e" />
      <path d="M 30 38 Q 30 22 42 22 Q 54 22 54 38" stroke="#0d9488" strokeWidth="6" fill="none" />
      <rect x="36" y="52" width="12" height="10" rx="2" fill="#fff" />
    </svg>
  ),
  "domain-service": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="42" cy="42" r="27" fill="none" stroke="#94A3B8" strokeWidth="5" strokeDasharray="7 7" />
      <text x="42" y="53" textAnchor="middle" fontSize="28" fontWeight="800" fill="#94A3B8">
        ?
      </text>
    </svg>
  ),
  factory: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="14" y="42" width="56" height="30" rx="4" fill="#C97A3D" />
      <polygon points="14,42 42,22 70,42" fill="#a85f28" />
      <rect x="37" y="54" width="10" height="18" fill="#fff" />
    </svg>
  ),
  "application-service": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="34" cy="28" r="14" fill="#3E6FB0" />
      <rect x="16" y="44" width="36" height="26" rx="10" fill="#3E6FB0" />
      <rect x="46" y="46" width="18" height="22" rx="3" fill="#fff" stroke="#3E6FB0" strokeWidth="2" />
      <line x1="50" y1="53" x2="60" y2="53" stroke="#3E6FB0" strokeWidth="2.4" strokeLinecap="round" />
      <line x1="50" y1="59" x2="60" y2="59" stroke="#3E6FB0" strokeWidth="2.4" strokeLinecap="round" />
    </svg>
  ),
  "domain-event": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <circle cx="42" cy="42" r="26" fill="none" stroke="#94A3B8" strokeWidth="5" strokeDasharray="7 7" />
      <circle cx="42" cy="42" r="9" fill="#94A3B8" />
    </svg>
  ),
  cqrs: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="10" y="12" width="26" height="26" rx="6" fill="#3E6FB0" />
      <rect x="48" y="12" width="26" height="26" rx="6" fill="#5B8F4E" />
      <ellipse cx="42" cy="62" rx="26" ry="9" fill="#7A5CA8" />
      <rect x="16" y="54" width="52" height="16" fill="#7A5CA8" />
      <ellipse cx="42" cy="70" rx="26" ry="9" fill="#654389" />
    </svg>
  ),
  "optimistic-locking": (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <rect x="18" y="36" width="44" height="34" rx="8" fill="#D97706" />
      <path d="M 30 36 a 10 10 0 0 1 20 0 v6 h-20 z" fill="none" stroke="#D97706" strokeWidth="4" />
      <circle cx="62" cy="24" r="14" fill="#fff" stroke="#D97706" strokeWidth="3" />
      <text x="62" y="29" textAnchor="middle" fontSize="11" fontWeight="800" fill="#D97706">
        v1
      </text>
    </svg>
  ),
  specification: (
    <svg viewBox="0 0 84 84" role="img" aria-hidden="true">
      <polygon points="16,18 68,18 46,48 38,48" fill="#94A3B8" />
      <rect x="38" y="48" width="8" height="22" fill="#94A3B8" />
    </svg>
  ),
};

export function TocGrid() {
  const shouldReduceMotion = useReducedMotion();
  const variants = shouldReduceMotion ? cardVariantsReduced : cardVariants;
  const sections = groupConceptPages(CONCEPT_PAGES);

  return (
    <div className="mt-9 flex flex-col gap-10">
      {sections.map((section) => (
        <div key={section.group}>
          <h3 className="mb-4 text-sm font-bold tracking-wide text-muted-foreground">
            {section.group}
          </h3>
          <motion.div
            className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-10%" }}
            variants={containerVariants}
          >
            {section.pages.map((page) => (
              <motion.div
                key={page.slug}
                variants={variants}
                whileHover={shouldReduceMotion ? undefined : { scale: 1.02 }}
                whileTap={shouldReduceMotion ? undefined : { scale: 0.97 }}
                transition={{ type: "spring", stiffness: 400, damping: 24 }}
              >
                <Card className="h-full overflow-visible rounded-3xl shadow-clay ring-1 ring-border">
                  <CardContent className="px-5">
                    <Link
                      href={pageHref(page)}
                      className="flex flex-col items-center gap-2.5 text-center text-foreground focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50 rounded-2xl"
                    >
                      <span className="flex h-20 w-20 items-center justify-center" aria-hidden="true">
                        {ICONS[page.slug]}
                      </span>
                      <span className="text-xs font-bold tracking-wide text-muted-foreground">
                        {page.num}
                      </span>
                      <span className="text-lg font-extrabold">{page.title}</span>
                      <span className="text-sm text-muted-foreground">{page.description}</span>
                    </Link>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </motion.div>
        </div>
      ))}
    </div>
  );
}
