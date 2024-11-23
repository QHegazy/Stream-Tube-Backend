import { EmailService } from "./emailService";

const emailService = new EmailService();

const emailOptions = {
  to: "asmo6928@gmail.com",
  subject: "Test Email from TypeScript",
  text: "This is a test email.",
  html: "<b>This is a test email.</b>",
};

emailService.sendEmail(emailOptions).catch((err) => console.error(err));
