<template>
    <div class="px-5 py-4 flex flex-col gap-3">
        <div class="flex items-center gap-3">
            <span class="size-3 rounded-full shrink-0" :class="STATE_DOT[status.state]"></span>
            <div class="flex flex-col min-w-0">
                <span class="font-semibold">{{ status.name }}</span>
                <span v-if="status.hostname" class="text-xs text-zinc-500 truncate">{{ status.hostname }}</span>
            </div>
            <div class="flex-1"></div>
            <span v-if="uptime !== undefined" class="text-sm text-zinc-500 cursor-default" :title="uptimeTooltip">
                {{ formatUptime(uptime) }} uptime
            </span>
            <span class="rounded-full px-2 py-0.5 text-xs font-medium" :class="STATE_BADGE[status.state]">{{ STATE_LABEL[status.state] }}</span>
        </div>

        <div v-if="status.history.length" class="flex gap-0.5 h-6">
            <div
                v-for="(point, idx) in history"
                :key="idx"
                class="flex-1 rounded-sm"
                :class="point ? (point.success ? 'bg-green-500' : 'bg-red-500') : 'bg-zinc-200 dark:bg-zinc-800'"
                :title="point ? `${point.success ? 'Up' : 'Down'} · ${new Date(point.timestamp).toLocaleString()}` : 'No data'"
            ></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { EndpointStatus, HistoryPoint, State } from "@/types/status";
import { computed } from "vue";

const HISTORY_BARS = 50;

const STATE_LABEL: Record<State, string> = { up: "Operational", down: "Down", unknown: "Unknown" };
const STATE_BADGE: Record<State, string> = {
    up: "bg-green-100 text-green-800 dark:bg-green-950 dark:text-green-300",
    down: "bg-red-100 text-red-800 dark:bg-red-950 dark:text-red-300",
    unknown: "bg-zinc-100 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-300",
};
const STATE_DOT: Record<State, string> = { up: "bg-green-500", down: "bg-red-500", unknown: "bg-zinc-400" };

const props = defineProps<{ status: EndpointStatus }>();

// Left-pad with empty bars so the history always renders as a fixed-width strip, newest on the right
const history = computed<(HistoryPoint | null)[]>(() => {
    const recent = props.status.history.slice(-HISTORY_BARS);
    return [...Array(HISTORY_BARS - recent.length).fill(null), ...recent];
});

const uptime = computed(() => props.status.uptime["30d"] ?? props.status.uptime["7d"] ?? props.status.uptime["24h"]);

const uptimeTooltip = computed(() =>
    ["24h", "7d", "30d"]
        .filter((d) => props.status.uptime[d] !== undefined)
        .map((d) => `${d}: ${formatUptime(props.status.uptime[d]!)}`)
        .join(" · "),
);

const formatUptime = (value: number) => `${(Math.floor(value * 10000) / 100).toFixed(2)}%`;
</script>
