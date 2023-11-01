package controller

import (
	"diary_api/helper"
	"diary_api/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddEntry(context *gin.Context) {
    var input model.Entry
    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := helper.CurrentUser(context)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.UserID = user.ID

    savedEntry, err := input.Save()

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusCreated, gin.H{"data": savedEntry})
}

func GetAllEntries(context *gin.Context) {
    user, err := helper.CurrentUser(context)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusOK, gin.H{"data": user.Entries})
}

func GetEntryByID(context *gin.Context) {
    user, err := helper.CurrentUser(context)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    id := context.Param("id")

    for _, entry := range user.Entries {
		if fmt.Sprint(entry.ID) == id {
            context.JSON(http.StatusOK, gin.H{"data": entry})
            return
        }
    }

    context.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
}

func UpdateEntry(context *gin.Context) {
    user, err := helper.CurrentUser(context)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    id := context.Param("id")
    var updatedEntry model.Entry

    if err := context.ShouldBindJSON(&updatedEntry); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for _, entry := range user.Entries {
        if fmt.Sprint(entry.ID) == id {
            // Update the entry with the new data
            if err := entry.Update(updatedEntry); err != nil {
                context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }

            context.JSON(http.StatusOK, gin.H{"data": entry})
            return
        }
    }

    context.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
}

func DeleteEntry(context *gin.Context) {
    user, err := helper.CurrentUser(context)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    id := context.Param("id")

    for i, entry := range user.Entries {
        if fmt.Sprint(entry.ID) == id {
            // Delete the entry
            if err := entry.Delete(); err != nil {
                context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }

            // Remove the entry from the user's list of entries
            user.Entries = append(user.Entries[:i], user.Entries[i+1:]...)

            context.JSON(http.StatusOK, gin.H{"message": "Entry deleted"})
            return
        }
    }

    context.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
}

