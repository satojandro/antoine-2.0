package cmd

import (
	"antoine-cli/internal/ui/views"
	"github.com/spf13/cobra"
)

var mentorCmd = &cobra.Command{
	Use:   "mentor",
	Short: "Interactive AI mentorship and guidance",
	Long: `Get personalized mentorship from Antoine. Whether you need project ideas,
code reviews, or strategic advice - Antoine learns from thousands of successful
hackathon projects to guide you to victory.`,
}

var mentorStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an interactive mentorship session",

	Run: func(cmd *cobra.Command, args []string) {
		view := views.NewMentorView(client)
		view.StartSession()
	},
}

var mentorFeedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Get feedback on a specific project",

	Run: func(cmd *cobra.Command, args []string) {
		options := &views.MentorOptions{
			ProjectURL: cmd.Flag("project-url").Value.String(),
			ProjectID:  cmd.Flag("project-id").Value.String(),
			Focus:      cmd.Flag("focus").Value.String(),
			Quick:      cmd.Flag("quick").Changed,
		}

		view := views.NewMentorView(client)
		view.ProvideFeedback(options)
	},
}

func init() {
	// Flags para feedback
	mentorFeedbackCmd.Flags().String("project-url", "", "GitHub repository URL")
	mentorFeedbackCmd.Flags().String("project-id", "", "Antoine project ID")
	mentorFeedbackCmd.Flags().StringSlice("focus", []string{}, "focus areas for feedback")
	mentorFeedbackCmd.Flags().Bool("quick", false, "quick feedback mode")

	mentorCmd.AddCommand(mentorStartCmd)
	mentorCmd.AddCommand(mentorFeedbackCmd)
}
