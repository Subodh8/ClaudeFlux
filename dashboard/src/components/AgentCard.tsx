'use client';

interface AgentCardProps {
    name: string;
    status: 'pending' | 'running' | 'complete' | 'failed' | 'blocked';
    tokensUsed?: number;
    costUSD?: number;
}

const statusColors = {
    pending: 'bg-cf-muted',
    running: 'bg-cf-primary animate-pulse',
    complete: 'bg-cf-success',
    failed: 'bg-cf-error',
    blocked: 'bg-cf-warning',
};

export function AgentCard({ name, status, tokensUsed = 0, costUSD = 0 }: AgentCardProps) {
    return (
        <div className="border border-cf-border rounded-lg p-4 bg-cf-surface">
            <div className="flex items-center justify-between mb-2">
                <span className="font-mono font-bold">{name}</span>
                <span className={`w-3 h-3 rounded-full ${statusColors[status]}`} />
            </div>
            <div className="text-sm text-cf-muted">
                <span>Tokens: {tokensUsed.toLocaleString()}</span>
                <span className="ml-4">Cost: ${costUSD.toFixed(4)}</span>
            </div>
        </div>
    );
}
