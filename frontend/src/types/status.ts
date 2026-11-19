export type State = "up" | "down" | "unknown";

export interface HistoryPoint {
    success: boolean;
    timestamp: string;
}

export interface EndpointStatus {
    name: string;
    group?: string;
    state: State;
    hostname?: string;
    lastChecked?: string;
    uptime: Record<string, number>;
    history: HistoryPoint[];
}
