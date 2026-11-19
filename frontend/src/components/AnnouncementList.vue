<template>
    <ul class="flex flex-col gap-3">
        <li
            v-for="(announcement, idx) in announcements"
            :key="idx"
            class="flex items-start gap-3 rounded-lg px-4 py-3"
            :class="archived ? 'bg-white ring-1 ring-zinc-200 dark:bg-zinc-900 dark:ring-zinc-800' : styleFor(announcement.type).class"
        >
            <component
                :is="styleFor(announcement.type).icon"
                class="size-5 shrink-0 mt-0.5"
                :class="archived && 'text-zinc-400'"
            />
            <div class="flex flex-col min-w-0">
                <span class="whitespace-pre-line break-words">{{ announcement.message }}</span>
                <time :datetime="announcement.timestamp" class="text-xs opacity-75">{{ new Date(announcement.timestamp).toLocaleString() }}</time>
            </div>
        </li>
    </ul>
</template>

<script setup lang="ts">
import type { Announcement, AnnouncementType } from "@/types/status";
import type { Component } from "vue";
import IconCircleCheck from "~icons/lucide/circle-check";
import IconCircleX from "~icons/lucide/circle-x";
import IconInfo from "~icons/lucide/info";
import IconMegaphone from "~icons/lucide/megaphone";
import IconTriangleAlert from "~icons/lucide/triangle-alert";

defineProps<{ announcements: Announcement[]; archived?: boolean }>();

const TYPE_STYLE: Record<AnnouncementType, { class: string; icon: Component }> = {
    outage: { class: "bg-red-100 text-red-900 dark:bg-red-950 dark:text-red-200", icon: IconCircleX },
    warning: { class: "bg-amber-100 text-amber-900 dark:bg-amber-950 dark:text-amber-200", icon: IconTriangleAlert },
    information: { class: "bg-blue-100 text-blue-900 dark:bg-blue-950 dark:text-blue-200", icon: IconInfo },
    operational: { class: "bg-green-100 text-green-900 dark:bg-green-950 dark:text-green-200", icon: IconCircleCheck },
    none: { class: "bg-zinc-100 text-zinc-800 dark:bg-zinc-800 dark:text-zinc-200", icon: IconMegaphone },
};
// Gatus may add announcement types in future, so fall back to the plain style
const styleFor = (type: AnnouncementType) => TYPE_STYLE[type] ?? TYPE_STYLE.none;
</script>
