'use client';

export function DAGVisualization() {
    // TODO: Integrate React Flow for interactive DAG rendering
    return (
        <div className="border border-cf-border rounded-lg p-8 bg-cf-surface">
            <p className="text-cf-muted text-center">
                DAG visualization will render here when a workflow is running.
            </p>
            <p className="text-cf-muted text-center mt-2 text-sm">
                Connect to the ClaudeFlux runtime SSE endpoint to see real-time updates.
            </p>
        </div>
    );
}
