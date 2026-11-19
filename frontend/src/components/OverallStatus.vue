<template>
    <div class="card shadow-sm text-white" :class="overall.class">
        <div class="card-body flex-row items-center gap-3 py-5">
            <component :is="overall.icon" class="size-7" />
            <span class="text-xl font-semibold">{{ overall.label }}</span>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { EndpointStatus } from "@/types/status";
import { computed } from "vue";
import IconCircleCheck from "~icons/lucide/circle-check";
import IconCircleX from "~icons/lucide/circle-x";
import IconCircleHelp from "~icons/lucide/circle-help";
import IconTriangleAlert from "~icons/lucide/triangle-alert";

const props = defineProps<{ statuses: EndpointStatus[] }>();

const overall = computed(() => {
    const down = props.statuses.filter((s) => s.state === "down").length;
    const unknown = props.statuses.filter((s) => s.state === "unknown").length;

    if (down === props.statuses.length) return { label: "Major Outage", class: "bg-error", icon: IconCircleX };
    if (down > 0) return { label: "Partial Outage", class: "bg-warning", icon: IconTriangleAlert };
    if (unknown > 0) return { label: "Some Systems Unknown", class: "bg-neutral", icon: IconCircleHelp };
    return { label: "All Systems Operational", class: "bg-success", icon: IconCircleCheck };
});
</script>
