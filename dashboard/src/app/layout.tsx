import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
    title: 'ClaudeFlux Dashboard',
    description: 'Real-time monitoring for ClaudeFlux multi-agent workflows',
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en" className="dark">
            <body className="bg-cf-bg text-cf-text min-h-screen">
                <nav className="border-b border-cf-border px-6 py-3 flex items-center gap-6">
                    <span className="font-bold text-lg">⚡ ClaudeFlux</span>
                    <a href="/" className="text-cf-muted hover:text-cf-text">DAG</a>
                    <a href="/logs" className="text-cf-muted hover:text-cf-text">Logs</a>
                    <a href="/costs" className="text-cf-muted hover:text-cf-text">Costs</a>
                    <a href="/approvals" className="text-cf-muted hover:text-cf-text">Approvals</a>
                </nav>
                <main className="p-6">{children}</main>
            </body>
        </html>
    );
}
