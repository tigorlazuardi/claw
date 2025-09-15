<script lang="ts">
  import { useQueryClient } from "@sveltestack/svelte-query";
  import { getSourceServiceClient } from "../../connectrpc";
  import { toDate } from "../../connectrpc/js_date";
  import CronExpressionParser from "cron-parser";

  const queryClient = useQueryClient();

  interface Props {
    value: string[];
  }

  let { value = $bindable([]) }: Props = $props();

  async function getNextRun(pattern: string) {
    return queryClient.fetchQuery(
      ["schedules", "nextRun", pattern],
      () =>
        getSourceServiceClient().then((client) =>
          client.getCronNextTime({
            cronExpression: pattern,
          }),
        ),
      {
        staleTime: 60 * 1000,
        cacheTime: 60 * 1000,
      },
    );
  }

  type ScheduleEntry = {
    pattern: string;
    nextRun?: Date;
    error?: string;
  };

  const schedules = $derived.by(async () => {
    const result: ScheduleEntry[] = [];
    for (const pattern of value) {
      const entry: ScheduleEntry = { pattern };
      try {
        CronExpressionParser.parse(pattern);
        const resp = await getNextRun(pattern);
        resp.nextTime;
        entry.nextRun = toDate(resp.nextTime);
      } catch (err) {
        if (err instanceof Error) {
          entry.error = err.message;
        } else {
          entry.error = "Invalid cron expression";
        }
      }
      result.push(entry);
    }
    return result;
  });

  type PatternOption = {
    label: string;
    pattern: string;
    description: string;
  };

  const patternDropdownList: PatternOption[] = [
    {
      label: "Every minute",
      pattern: "* * * * *",
      description: "Runs every minute",
    },
    {
      label: "Every minute (alt)",
      pattern: "@minutely",
      description: "Runs every minute",
    },
    {
      label: "Every 5 minutes",
      pattern: "*/5 * * * *",
      description: "Runs every 5 minutes",
    },
    {
      label: "Every 15 minutes",
      pattern: "*/15 * * * *",
      description: "Runs every 15 minutes",
    },
    {
      label: "Every 30 minutes",
      pattern: "*/30 * * * *",
      description: "Runs every 30 minutes",
    },
    {
      label: "Every hour",
      pattern: "0 * * * *",
      description: "Runs at the start of every hour",
    },
    {
      label: "Daily at midnight",
      pattern: "0 0 * * *",
      description: "Runs daily at midnight",
    },
    {
      label: "Weekly on Sunday at midnight",
      pattern: "0 0 * * 0",
      description: "Runs weekly on Sunday at midnight",
    },
    {
      label: "Monthly on the 1st at midnight",
      pattern: "0 0 1 * *",
      description: "Runs monthly on the 1st at midnight",
    },
  ];
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span>Schedule</span>
  </legend>
</fieldset>
