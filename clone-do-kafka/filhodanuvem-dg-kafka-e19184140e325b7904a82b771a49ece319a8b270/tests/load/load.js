import exec from 'k6/x/exec';

export default function () {
    console.log(exec.command("kafka-cli", ["-p", "-t", "orders", "-m", "payment1"]));
}