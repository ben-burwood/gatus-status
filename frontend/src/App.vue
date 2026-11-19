<template>
    <div class="min-h-screen w-full flex flex-col bg-zinc-50 text-zinc-900 dark:bg-zinc-950 dark:text-zinc-100">
        <main class="flex-1 w-full max-w-3xl mx-auto px-5 py-12 flex flex-col gap-6">
            <div class="text-center">
                <img v-if="logo" :src="logo" alt="" class="mx-auto mb-4 max-h-16 max-w-48 object-contain" />
                <h1 class="text-3xl font-bold" :class="{ invisible: !title }">{{ title || "Status" }}</h1>
                <p v-if="updatedAt" class="mt-1 text-sm text-zinc-500">Updated {{ updatedAt.toLocaleTimeString() }}</p>
            </div>

            <div
                v-if="error"
                role="alert"
                class="flex items-center gap-3 rounded-lg px-4 py-3 bg-amber-100 text-amber-900 dark:bg-amber-950 dark:text-amber-200"
            >
                <IconTriangleAlert class="size-5 shrink-0" />
                <span>Unable to load the current status. Retrying…</span>
            </div>

            <AnnouncementList v-if="activeAnnouncements.length" :announcements="activeAnnouncements" />

            <template v-if="statuses.length">
                <OverallStatus :statuses="statuses" />
                <ul class="rounded-lg bg-white shadow-sm ring-1 ring-zinc-200 divide-y divide-zinc-200 dark:bg-zinc-900 dark:ring-zinc-800 dark:divide-zinc-800">
                    <li v-for="status in statuses" :key="status.group + '/' + status.name">
                        <EndpointRow :status="status" />
                    </li>
                </ul>
            </template>
            <div v-else-if="loading" class="flex justify-center mt-10">
                <IconLoaderCircle class="size-8 animate-spin text-zinc-400" />
            </div>

            <section v-if="pastAnnouncements.length" class="mt-4 flex flex-col gap-3">
                <h2 class="text-lg font-semibold">Past Announcements</h2>
                <AnnouncementList :announcements="pastAnnouncements" archived />
            </section>
        </main>

        <footer class="py-6 flex items-center justify-center gap-3 text-sm text-zinc-500">
            <span>
                Powered by
                <a href="https://github.com/TwiN/gatus" target="_blank" rel="noopener" class="font-medium hover:underline">Gatus</a>
                <span class="mx-1">·</span>
                <a href="https://github.com/ben-burwood/gatus-status" target="_blank" rel="noopener" class="font-medium hover:underline">gatus-status</a>
            </span>
            <ThemeToggle />
        </footer>
    </div>
</template>

<script setup lang="ts">
import AnnouncementList from "@/components/AnnouncementList.vue";
import EndpointRow from "@/components/EndpointRow.vue";
import OverallStatus from "@/components/OverallStatus.vue";
import ThemeToggle from "@/components/ThemeToggle.vue";
import { SERVER_URL } from "@/main";
import type { Announcement, EndpointStatus } from "@/types/status";
import { computed, onMounted, onUnmounted, ref } from "vue";
import IconLoaderCircle from "~icons/lucide/loader-circle";
import IconTriangleAlert from "~icons/lucide/triangle-alert";

const REFRESH_INTERVAL_MS = 60 * 1000;

// Hidden until the configured title has loaded, to avoid flashing a placeholder
const title = ref("");
const statuses = ref<EndpointStatus[]>([]);
const loading = ref(true);
const error = ref(false);
const updatedAt = ref<Date | null>(null);
const logo = ref("");
const showAnnouncements = ref(false);
const announcements = ref<Announcement[]>([]);
let refreshHandler: ReturnType<typeof setInterval> | null = null;

const activeAnnouncements = computed(() => announcements.value.filter((a) => !a.archived));
const pastAnnouncements = computed(() => announcements.value.filter((a) => a.archived));

const loadConfig = async () => {
    try {
        const res = await fetch(`${SERVER_URL}/config`);
        const config: { title: string; logo: string; announcements: boolean } = await res.json();
        title.value = config.title;
        document.title = config.title;
        logo.value = config.logo;
        showAnnouncements.value = config.announcements;
    } catch (e) {
        console.error("Failed to fetch config:", e);
        title.value = "Status";
    }
};

// Failures keep the last announcements shown, as they're secondary to the statuses
const loadAnnouncements = async () => {
    if (!showAnnouncements.value) return;
    try {
        const res = await fetch(`${SERVER_URL}/announcements`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        announcements.value = await res.json();
    } catch (e) {
        console.error("Failed to fetch announcements:", e);
    }
};

const refresh = () => {
    loadStatuses();
    loadAnnouncements();
};

const loadStatuses = async () => {
    try {
        const res = await fetch(`${SERVER_URL}/statuses`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        statuses.value = await res.json();
        updatedAt.value = new Date();
        error.value = false;
    } catch (e) {
        console.error("Failed to fetch statuses:", e);
        error.value = true;
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    loadConfig().then(loadAnnouncements);
    loadStatuses();
    refreshHandler = setInterval(refresh, REFRESH_INTERVAL_MS);
});

onUnmounted(() => {
    if (refreshHandler) clearInterval(refreshHandler);
});
</script>
