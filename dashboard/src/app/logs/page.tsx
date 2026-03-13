import { LogStream } from '../../components/LogStream';

export default function LogsPage() {
    return (
        <div>
            <h1 className="text-2xl font-bold mb-6">Agent Logs</h1>
            <LogStream />
        </div>
    );
}
