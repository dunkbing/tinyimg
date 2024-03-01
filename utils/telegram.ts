import { TelegramBot } from "https://deno.land/x/telegram_bot_api@0.4.0/mod.ts";
import config from "@/utils/config.ts";

let _bot: TelegramBot;

function getBot() {
  try {
    if (!_bot) {
      _bot = new TelegramBot(config.teleBotToken as string);
    }
    return _bot;
  } catch (e) {
    console.error("Error creating telegram bot", e);
    return null;
  }
}

export function sendMessage(text: string) {
  const bot = getBot();
  return bot?.sendMessage({
    chat_id: config.teleChatID as string,
    text,
  }).catch((err) => console.error("Send tele message error:", err));
}
