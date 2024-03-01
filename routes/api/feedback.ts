import { FreshContext } from "$fresh/server.ts";
import { sendMessage } from "@/utils/telegram.ts";

type FeedbackData = {
  subject: string;
  message: string;
};

export const handler = {
  async POST(req: Request, _ctx: FreshContext): Promise<Response> {
    const { subject, message } = (await req.json()) as FeedbackData;
    const text = [
      "TinyIMG Feedback",
      "---------------",
      subject,
      "---------------",
      message,
    ].join("\n");
    void sendMessage(text);

    return new Response(JSON.stringify({ message: "ok" }));
  },
};
