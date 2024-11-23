import nodemailer from "nodemailer";
import dotenv from "dotenv";

dotenv.config();

interface EmailOptions {
  to: string;
  subject: string;
  text?: string;
  html?: string;
}

export class EmailService {
  private transporter: nodemailer.Transporter;

  constructor() {
    this.transporter = nodemailer.createTransport({
      service: "gmail", 
      auth: {
        user: process.env.GMAIL_USER, 
        pass: process.env.GMAIL_PASS, 
      },
    });
  }

  public async sendEmail(options: EmailOptions): Promise<void> {
    try {
      const mailOptions = {
        from: process.env.GMAIL_USER, 
        to: options.to, 
        subject: options.subject, 
        text: options.text || "",
        html: options.html || "",
      };

      await this.transporter.sendMail(mailOptions);
      console.log(`Email sent to ${options.to}`);
    } catch (error: unknown) {
      if (error instanceof Error) {
        console.error(`Failed to send email: ${error.message}`);
      } else {
        console.error("An unknown error occurred");
      }
      throw error; // Optional: rethrow the error if needed
    }
  }
}
