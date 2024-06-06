package main

import (
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

const SOLVEDAPI string = "https://solved.ac/api/v3"

type UserInfo struct {
	Handle                  string    `json:"handle"`
	Bio                     string    `json:"bio"`
	BadgeID                 string    `json:"badgeId"`
	BackgroundID            string    `json:"backgroundId"`
	ProfileImageURL         string    `json:"profileImageUrl"`
	SolvedCount             int       `json:"solvedCount"`
	VoteCount               int       `json:"voteCount"`
	Class                   int       `json:"class"`
	ClassDecoration         string    `json:"classDecoration"`
	RivalCount              int       `json:"rivalCount"`
	ReverseRivalCount       int       `json:"reverseRivalCount"`
	Tier                    int       `json:"tier"`
	Rating                  int       `json:"rating"`
	RatingByProblemsSum     int       `json:"ratingByProblemsSum"`
	RatingByClass           int       `json:"ratingByClass"`
	RatingBySolvedCount     int       `json:"ratingBySolvedCount"`
	RatingByVoteCount       int       `json:"ratingByVoteCount"`
	ArenaTier               int       `json:"arenaTier"`
	ArenaRating             int       `json:"arenaRating"`
	ArenaMaxTier            int       `json:"arenaMaxTier"`
	ArenaMaxRating          int       `json:"arenaMaxRating"`
	ArenaCompetedRoundCount int       `json:"arenaCompetedRoundCount"`
	MaxStreak               int       `json:"maxStreak"`
	Coins                   int       `json:"coins"`
	Stardusts               int       `json:"stardusts"`
	JoinedAt                time.Time `json:"joinedAt"`
	BannedUntil             time.Time `json:"bannedUntil"`
	ProUntil                time.Time `json:"proUntil"`
	Rank                    int       `json:"rank"`
	IsRival                 bool      `json:"isRival"`
	IsReverseRival          bool      `json:"isReverseRival"`
	Blocked                 bool      `json:"blocked"`
	ReverseBlocked          bool      `json:"reverseBlocked"`
}

func GetUserInfo(handle string) (*UserInfo, error) {
	client := req.C()
	userInfo := UserInfo{}
	url := fmt.Sprintf("%s/user/show", SOLVEDAPI)
	resp, err := client.R().
		SetQueryParam("handle", handle).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetSuccessResult(&userInfo).
		Get(url)

	if err != nil || resp.IsErrorState() {
		if err == nil {
			err = fmt.Errorf("status code: %d, error: %v", resp.StatusCode, resp.String())
		}
		return nil, err
	}

	return &userInfo, nil
}
