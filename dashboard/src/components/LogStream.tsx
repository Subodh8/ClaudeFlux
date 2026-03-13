'use client';

export function LogStream() {
    // TODO: Connect to SSE endpoint for real-time log streaming
    return (
        <div className="border border-cf-border rounded-lg bg-cf-surface font-mono text-sm">
            <div className="px-4 py-2 border-b border-cf-border flex items-center gap-4">
                <span className="text-cf-muted">Filter by agent:</span>
                <select className="bg-cf-bg border border-cf-border rounded px-2 py-1 text-cf-text">
                    <option value="">All agents</option>
                </select>
            </div>
            <div className="p-4 min-h-[400px]">
                <p className="text-cf-muted">Waiting for log events...</p>
            </div>
        </div>
    );
}
